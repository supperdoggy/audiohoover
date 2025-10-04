//go:build integration
// +build integration

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"github.com/supperdoggy/audiohoover/pkg/service"
	"go.uber.org/zap"
)

func TestEndToEndWorkflow(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	playlistsDir := filepath.Join(tmpDir, "playlists")
	outputDir := filepath.Join(tmpDir, "output")

	// Create playlists structure
	playlists := map[string][]string{
		"rock": {"song1.mp3", "song2.flac", "cover.jpg"},
		"jazz": {"track1.wav", "track2.mp3"},
		"pop":  {"hit1.m4a", "hit2.aac", "readme.txt"},
	}

	for playlistName, files := range playlists {
		playlistPath := filepath.Join(playlistsDir, playlistName)
		if err := os.MkdirAll(playlistPath, 0755); err != nil {
			t.Fatalf("failed to create playlist %s: %v", playlistName, err)
		}

		for _, filename := range files {
			filePath := filepath.Join(playlistPath, filename)
			content := []byte("test content for " + filename)
			if err := os.WriteFile(filePath, content, 0644); err != nil {
				t.Fatalf("failed to create file %s: %v", filename, err)
			}
		}
	}

	// Set up environment
	os.Setenv("PLAYLISTS_ROOT", playlistsDir)
	os.Setenv("DESTINATION_FOLDER", outputDir)
	os.Setenv("AUDIO_EXTENSIONS", ".mp3,.flac,.wav,.m4a,.aac")
	defer os.Clearenv()

	// Create config and service
	ctx := context.Background()
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	log, _ := zap.NewDevelopment()
	svc := service.NewService(log, cfg)

	// Run the service
	stats, err := svc.Run()
	if err != nil {
		t.Fatalf("service run failed: %v", err)
	}

	// Validate results
	expectedCopied := 7 // 7 audio files
	if stats.FilesCopied != expectedCopied {
		t.Errorf("expected %d files copied, got %d", expectedCopied, stats.FilesCopied)
	}

	expectedSkipped := 2 // cover.jpg and readme.txt
	if stats.FilesSkipped != expectedSkipped {
		t.Errorf("expected %d files skipped, got %d", expectedSkipped, stats.FilesSkipped)
	}

	// Verify all audio files exist in output
	audioFiles := []string{"song1.mp3", "song2.flac", "track1.wav", "track2.mp3", "hit1.m4a", "hit2.aac"}
	for _, filename := range audioFiles {
		path := filepath.Join(outputDir, filename)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s not found in output: %v", filename, err)
		}
	}

	// Verify non-audio files don't exist in output
	nonAudioFiles := []string{"cover.jpg", "readme.txt"}
	for _, filename := range nonAudioFiles {
		path := filepath.Join(outputDir, filename)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("non-audio file %s should not be in output", filename)
		}
	}

	// Run again to test duplicate detection
	stats2, err := svc.Run()
	if err != nil {
		t.Fatalf("second service run failed: %v", err)
	}

	if stats2.FilesCopied != 0 {
		t.Errorf("second run should copy 0 files (all duplicates), got %d", stats2.FilesCopied)
	}

	if stats2.FilesSkipped != 9 { // 7 audio + 2 non-audio
		t.Errorf("second run should skip 9 files, got %d", stats2.FilesSkipped)
	}
}

func TestDryRunMode(t *testing.T) {
	tmpDir := t.TempDir()
	playlistsDir := filepath.Join(tmpDir, "playlists")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test structure
	playlist := filepath.Join(playlistsDir, "test")
	if err := os.MkdirAll(playlist, 0755); err != nil {
		t.Fatalf("failed to create playlist: %v", err)
	}

	testFile := filepath.Join(playlist, "song.mp3")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set up dry-run config
	os.Setenv("PLAYLISTS_ROOT", playlistsDir)
	os.Setenv("DESTINATION_FOLDER", outputDir)
	os.Setenv("DRY_RUN", "true")
	defer os.Clearenv()

	ctx := context.Background()
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	log, _ := zap.NewDevelopment()
	svc := service.NewService(log, cfg)

	stats, err := svc.Run()
	if err != nil {
		t.Fatalf("service run failed: %v", err)
	}

	// Should report 1 file "copied"
	if stats.FilesCopied != 1 {
		t.Errorf("expected 1 file copied in dry-run, got %d", stats.FilesCopied)
	}

	// But file should not actually exist
	outputFile := filepath.Join(outputDir, "song.mp3")
	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		t.Error("file should not exist in dry-run mode")
	}
}
