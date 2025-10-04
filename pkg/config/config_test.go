package config

import (
	"context"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid config with defaults",
			envVars: map[string]string{
				"PLAYLISTS_ROOT": tmpDir,
			},
			wantErr: false,
		},
		{
			name: "valid config with all options",
			envVars: map[string]string{
				"PLAYLISTS_ROOT":     tmpDir,
				"DESTINATION_FOLDER": tmpDir,
				"DRY_RUN":            "true",
				"AUDIO_EXTENSIONS":   ".mp3,.wav",
				"VERBOSE":            "true",
			},
			wantErr: false,
		},
		{
			name: "playlists root does not exist",
			envVars: map[string]string{
				"PLAYLISTS_ROOT": "/nonexistent/path",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			ctx := context.Background()
			cfg, err := NewConfig(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("NewConfig() returned nil config without error")
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				PlaylistsRoot:     tmpDir,
				DestinationFolder: tmpDir,
				AudioExtensions:   ".mp3,.wav",
			},
			wantErr: false,
		},
		{
			name: "empty playlists root",
			config: &Config{
				PlaylistsRoot:     "",
				DestinationFolder: tmpDir,
				AudioExtensions:   ".mp3",
			},
			wantErr: true,
		},
		{
			name: "empty destination folder",
			config: &Config{
				PlaylistsRoot:     tmpDir,
				DestinationFolder: "",
				AudioExtensions:   ".mp3",
			},
			wantErr: true,
		},
		{
			name: "playlists root does not exist",
			config: &Config{
				PlaylistsRoot:     "/nonexistent/path",
				DestinationFolder: tmpDir,
				AudioExtensions:   ".mp3",
			},
			wantErr: true,
		},
		{
			name: "empty audio extensions",
			config: &Config{
				PlaylistsRoot:     tmpDir,
				DestinationFolder: tmpDir,
				AudioExtensions:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetAudioExtensions(t *testing.T) {
	tests := []struct {
		name       string
		extensions string
		want       []string
	}{
		{
			name:       "single extension with dot",
			extensions: ".mp3",
			want:       []string{".mp3"},
		},
		{
			name:       "multiple extensions with dots",
			extensions: ".mp3,.wav,.flac",
			want:       []string{".mp3", ".wav", ".flac"},
		},
		{
			name:       "extensions without dots",
			extensions: "mp3,wav",
			want:       []string{".mp3", ".wav"},
		},
		{
			name:       "extensions with spaces",
			extensions: ".mp3, .wav , .flac",
			want:       []string{".mp3", ".wav", ".flac"},
		},
		{
			name:       "mixed case extensions",
			extensions: ".MP3,.WaV,.FLAC",
			want:       []string{".mp3", ".wav", ".flac"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AudioExtensions: tt.extensions,
			}
			got := cfg.GetAudioExtensions()
			if len(got) != len(tt.want) {
				t.Errorf("GetAudioExtensions() returned %d extensions, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetAudioExtensions()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestConfig_IsAudioFile(t *testing.T) {
	cfg := &Config{
		AudioExtensions: ".mp3,.wav,.flac",
	}

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "mp3 file",
			filename: "song.mp3",
			want:     true,
		},
		{
			name:     "wav file",
			filename: "audio.wav",
			want:     true,
		},
		{
			name:     "flac file uppercase",
			filename: "music.FLAC",
			want:     true,
		},
		{
			name:     "txt file",
			filename: "readme.txt",
			want:     false,
		},
		{
			name:     "no extension",
			filename: "noextension",
			want:     false,
		},
		{
			name:     "different audio format",
			filename: "song.ogg",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.IsAudioFile(tt.filename)
			if got != tt.want {
				t.Errorf("IsAudioFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
