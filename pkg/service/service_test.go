package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"go.uber.org/zap"
)

func setupTestEnvironment(t *testing.T) (string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	playlistsDir := filepath.Join(tmpDir, "playlists")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(playlistsDir, 0755); err != nil {
		t.Fatalf("failed to create playlists dir: %v", err)
	}

	return playlistsDir, outputDir
}

func createTestFile(t *testing.T, dir, filename, content string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", filename, err)
	}
}

func TestService_Run(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T, playlistsDir, outputDir string) *config.Config
		validateFunc func(t *testing.T, stats *Stats, outputDir string)
		wantErr      bool
	}{
		{
			name: "successful processing",
			setupFunc: func(t *testing.T, playlistsDir, outputDir string) *config.Config {
				// Create a playlist with audio files
				playlist1 := filepath.Join(playlistsDir, "playlist1")
				if err := os.Mkdir(playlist1, 0755); err != nil {
					t.Fatalf("failed to create playlist1: %v", err)
				}
				createTestFile(t, playlist1, "song1.mp3", "audio content 1")
				createTestFile(t, playlist1, "song2.wav", "audio content 2")

				return &config.Config{
					PlaylistsRoot:     playlistsDir,
					DestinationFolder: outputDir,
					AudioExtensions:   ".mp3,.wav,.flac",
					DryRun:            false,
					Verbose:           false,
				}
			},
			validateFunc: func(t *testing.T, stats *Stats, outputDir string) {
				if stats.FilesCopied != 2 {
					t.Errorf("expected 2 files copied, got %d", stats.FilesCopied)
				}
				if stats.FilesSkipped != 0 {
					t.Errorf("expected 0 files skipped, got %d", stats.FilesSkipped)
				}
				// Verify files exist
				if _, err := os.Stat(filepath.Join(outputDir, "song1.mp3")); err != nil {
					t.Errorf("song1.mp3 not found in output: %v", err)
				}
				if _, err := os.Stat(filepath.Join(outputDir, "song2.wav")); err != nil {
					t.Errorf("song2.wav not found in output: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name: "skip duplicates",
			setupFunc: func(t *testing.T, playlistsDir, outputDir string) *config.Config {
				// Create output dir with existing file
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					t.Fatalf("failed to create output dir: %v", err)
				}
				createTestFile(t, outputDir, "duplicate.mp3", "existing content")

				// Create playlist with duplicate file
				playlist1 := filepath.Join(playlistsDir, "playlist1")
				if err := os.Mkdir(playlist1, 0755); err != nil {
					t.Fatalf("failed to create playlist1: %v", err)
				}
				createTestFile(t, playlist1, "duplicate.mp3", "new content")
				createTestFile(t, playlist1, "unique.mp3", "unique content")

				return &config.Config{
					PlaylistsRoot:     playlistsDir,
					DestinationFolder: outputDir,
					AudioExtensions:   ".mp3",
					DryRun:            false,
					Verbose:           false,
				}
			},
			validateFunc: func(t *testing.T, stats *Stats, outputDir string) {
				if stats.FilesCopied != 1 {
					t.Errorf("expected 1 file copied, got %d", stats.FilesCopied)
				}
				if stats.FilesSkipped != 1 {
					t.Errorf("expected 1 file skipped, got %d", stats.FilesSkipped)
				}
				// Verify original content preserved
				content, err := os.ReadFile(filepath.Join(outputDir, "duplicate.mp3"))
				if err != nil {
					t.Fatalf("failed to read duplicate.mp3: %v", err)
				}
				if string(content) != "existing content" {
					t.Errorf("duplicate.mp3 was overwritten")
				}
			},
			wantErr: false,
		},
		{
			name: "filter non-audio files",
			setupFunc: func(t *testing.T, playlistsDir, outputDir string) *config.Config {
				playlist1 := filepath.Join(playlistsDir, "playlist1")
				if err := os.Mkdir(playlist1, 0755); err != nil {
					t.Fatalf("failed to create playlist1: %v", err)
				}
				createTestFile(t, playlist1, "song.mp3", "audio")
				createTestFile(t, playlist1, "readme.txt", "text")
				createTestFile(t, playlist1, "image.jpg", "image")

				return &config.Config{
					PlaylistsRoot:     playlistsDir,
					DestinationFolder: outputDir,
					AudioExtensions:   ".mp3,.wav",
					DryRun:            false,
					Verbose:           false,
				}
			},
			validateFunc: func(t *testing.T, stats *Stats, outputDir string) {
				if stats.FilesCopied != 1 {
					t.Errorf("expected 1 file copied, got %d", stats.FilesCopied)
				}
				if stats.FilesSkipped != 2 {
					t.Errorf("expected 2 files skipped, got %d", stats.FilesSkipped)
				}
				if stats.FilesProcessed != 3 {
					t.Errorf("expected 3 files processed, got %d", stats.FilesProcessed)
				}
			},
			wantErr: false,
		},
		{
			name: "dry run mode",
			setupFunc: func(t *testing.T, playlistsDir, outputDir string) *config.Config {
				playlist1 := filepath.Join(playlistsDir, "playlist1")
				if err := os.Mkdir(playlist1, 0755); err != nil {
					t.Fatalf("failed to create playlist1: %v", err)
				}
				createTestFile(t, playlist1, "song.mp3", "audio")

				return &config.Config{
					PlaylistsRoot:     playlistsDir,
					DestinationFolder: outputDir,
					AudioExtensions:   ".mp3",
					DryRun:            true,
					Verbose:           false,
				}
			},
			validateFunc: func(t *testing.T, stats *Stats, outputDir string) {
				if stats.FilesCopied != 1 {
					t.Errorf("expected 1 file 'copied' (dry-run), got %d", stats.FilesCopied)
				}
				// Verify file was NOT actually copied
				if _, err := os.Stat(filepath.Join(outputDir, "song.mp3")); !os.IsNotExist(err) {
					t.Errorf("file should not exist in dry-run mode")
				}
			},
			wantErr: false,
		},
		{
			name: "multiple playlists",
			setupFunc: func(t *testing.T, playlistsDir, outputDir string) *config.Config {
				playlist1 := filepath.Join(playlistsDir, "playlist1")
				playlist2 := filepath.Join(playlistsDir, "playlist2")
				if err := os.Mkdir(playlist1, 0755); err != nil {
					t.Fatalf("failed to create playlist1: %v", err)
				}
				if err := os.Mkdir(playlist2, 0755); err != nil {
					t.Fatalf("failed to create playlist2: %v", err)
				}
				createTestFile(t, playlist1, "song1.mp3", "audio1")
				createTestFile(t, playlist2, "song2.mp3", "audio2")

				return &config.Config{
					PlaylistsRoot:     playlistsDir,
					DestinationFolder: outputDir,
					AudioExtensions:   ".mp3",
					DryRun:            false,
					Verbose:           false,
				}
			},
			validateFunc: func(t *testing.T, stats *Stats, outputDir string) {
				if stats.FilesCopied != 2 {
					t.Errorf("expected 2 files copied, got %d", stats.FilesCopied)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlistsDir, outputDir := setupTestEnvironment(t)
			cfg := tt.setupFunc(t, playlistsDir, outputDir)

			log, _ := zap.NewDevelopment()
			svc := NewService(log, cfg)

			stats, err := svc.Run()

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, stats, outputDir)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	log, _ := zap.NewDevelopment()
	cfg := &config.Config{
		PlaylistsRoot:     "/tmp/playlists",
		DestinationFolder: "/tmp/output",
		AudioExtensions:   ".mp3",
	}

	svc := NewService(log, cfg)

	if svc == nil {
		t.Fatal("NewService returned nil")
	}
	if svc.log != log {
		t.Error("logger not set correctly")
	}
	if svc.cfg != cfg {
		t.Error("config not set correctly")
	}
	if svc.copied == nil {
		t.Error("copied map not initialized")
	}
}
