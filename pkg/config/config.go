package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

var (
	// ErrInvalidConfig is returned when configuration validation fails
	ErrInvalidConfig = errors.New("invalid configuration")
)

// Config holds the application configuration
type Config struct {
	// PlaylistsRoot is the root directory containing playlists with audio files
	PlaylistsRoot string `env:"PLAYLISTS_ROOT, default=./playlists"`

	// DestinationFolder is where audio files will be copied to
	DestinationFolder string `env:"DESTINATION_FOLDER, default=./output"`

	// DryRun enables dry-run mode where no files are actually copied
	DryRun bool `env:"DRY_RUN, default=false"`

	// AudioExtensions is a comma-separated list of audio file extensions to process
	AudioExtensions string `env:"AUDIO_EXTENSIONS, default=.mp3,.flac,.wav,.m4a,.aac,.ogg,.wma"`

	// Verbose enables verbose logging
	Verbose bool `env:"VERBOSE, default=false"`
}

// NewConfig creates a new configuration from environment variables and validates it
func NewConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.PlaylistsRoot == "" {
		return fmt.Errorf("%w: playlists root cannot be empty", ErrInvalidConfig)
	}

	if c.DestinationFolder == "" {
		return fmt.Errorf("%w: destination folder cannot be empty", ErrInvalidConfig)
	}

	// Check if playlists root exists
	if info, err := os.Stat(c.PlaylistsRoot); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: playlists root does not exist: %s", ErrInvalidConfig, c.PlaylistsRoot)
		}
		return fmt.Errorf("failed to stat playlists root: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("%w: playlists root is not a directory: %s", ErrInvalidConfig, c.PlaylistsRoot)
	}

	// Validate audio extensions
	if c.AudioExtensions == "" {
		return fmt.Errorf("%w: audio extensions cannot be empty", ErrInvalidConfig)
	}

	return nil
}

// GetAudioExtensions returns the list of audio file extensions
func (c *Config) GetAudioExtensions() []string {
	parts := strings.Split(c.AudioExtensions, ",")
	extensions := make([]string, 0, len(parts))
	for _, ext := range parts {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			// Ensure extension starts with a dot
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extensions = append(extensions, strings.ToLower(ext))
		}
	}
	return extensions
}

// IsAudioFile checks if the given filename has an audio extension
func (c *Config) IsAudioFile(filename string) bool {
	filename = strings.ToLower(filename)
	for _, ext := range c.GetAudioExtensions() {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
