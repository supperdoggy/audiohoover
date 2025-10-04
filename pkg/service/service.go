package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"github.com/supperdoggy/audiohoover/pkg/utils"
	"go.uber.org/zap"
)

// Stats holds statistics about the processing run
type Stats struct {
	FilesProcessed int
	FilesCopied    int
	FilesSkipped   int
	Errors         []error
}

// Service handles the audio file processing logic
type Service struct {
	log    *zap.Logger
	cfg    *config.Config
	copied map[string]struct{} // tracks copied files to avoid duplicates
}

// NewService creates a new Service instance
func NewService(log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		log:    log,
		cfg:    cfg,
		copied: make(map[string]struct{}),
	}
}

// Run executes the main application logic
func (s *Service) Run() (*Stats, error) {
	stats := &Stats{
		Errors: make([]error, 0),
	}

	// Create destination folder if it doesn't exist
	if err := s.ensureDestinationExists(); err != nil {
		return stats, fmt.Errorf("failed to ensure destination exists: %w", err)
	}

	// Load existing files from destination
	if err := s.loadExistingFiles(); err != nil {
		return stats, fmt.Errorf("failed to load existing files: %w", err)
	}

	// Process playlists
	if err := s.processPlaylists(stats); err != nil {
		return stats, fmt.Errorf("failed to process playlists: %w", err)
	}

	return stats, nil
}

// ensureDestinationExists creates the destination directory if it doesn't exist
func (s *Service) ensureDestinationExists() error {
	if s.cfg.DryRun {
		s.log.Info("dry-run: would create destination folder", zap.String("path", s.cfg.DestinationFolder))
		return nil
	}

	if err := os.MkdirAll(s.cfg.DestinationFolder, 0755); err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	return nil
}

// loadExistingFiles reads the destination folder and populates the copied map
func (s *Service) loadExistingFiles() error {
	// Check if destination exists
	info, err := os.Stat(s.cfg.DestinationFolder)
	if err != nil {
		if os.IsNotExist(err) {
			// Destination doesn't exist yet, that's fine
			s.log.Info("destination folder does not exist yet", zap.String("path", s.cfg.DestinationFolder))
			return nil
		}
		return fmt.Errorf("failed to stat destination folder: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("destination path is not a directory: %s", s.cfg.DestinationFolder)
	}

	destinationFolder, err := os.Open(s.cfg.DestinationFolder)
	if err != nil {
		return fmt.Errorf("failed to open destination folder: %w", err)
	}
	defer destinationFolder.Close()

	files, err := destinationFolder.Readdir(0)
	if err != nil {
		return fmt.Errorf("failed to read destination folder: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			s.copied[file.Name()] = struct{}{}
		}
	}

	s.log.Info("loaded existing files", zap.Int("count", len(s.copied)))
	return nil
}

// processPlaylists processes all playlist folders
func (s *Service) processPlaylists(stats *Stats) error {
	playlistsFolder, err := os.Open(s.cfg.PlaylistsRoot)
	if err != nil {
		return fmt.Errorf("failed to open playlists root: %w", err)
	}
	defer playlistsFolder.Close()

	playlists, err := playlistsFolder.Readdir(0)
	if err != nil {
		return fmt.Errorf("failed to read playlists root: %w", err)
	}

	s.log.Info("found playlists", zap.Int("count", len(playlists)))

	for _, playlist := range playlists {
		if !playlist.IsDir() {
			s.log.Debug("skipping non-directory", zap.String("name", playlist.Name()))
			continue
		}

		playlistPath := filepath.Join(s.cfg.PlaylistsRoot, playlist.Name())
		s.log.Info("processing playlist", zap.String("name", playlist.Name()))

		if err := s.processPlaylist(playlistPath, stats); err != nil {
			s.log.Error("failed to process playlist",
				zap.String("playlist", playlist.Name()),
				zap.Error(err))
			stats.Errors = append(stats.Errors, fmt.Errorf("playlist %s: %w", playlist.Name(), err))
			continue
		}
	}

	return nil
}

// processPlaylist processes a single playlist folder
func (s *Service) processPlaylist(playlistPath string, stats *Stats) error {
	playlistFolder, err := os.Open(playlistPath)
	if err != nil {
		return fmt.Errorf("failed to open playlist folder: %w", err)
	}
	defer playlistFolder.Close()

	files, err := playlistFolder.Readdir(0)
	if err != nil {
		return fmt.Errorf("failed to read playlist folder: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		stats.FilesProcessed++

		// Check if it's an audio file
		if !s.cfg.IsAudioFile(file.Name()) {
			if s.cfg.Verbose {
				s.log.Debug("skipping non-audio file", zap.String("file", file.Name()))
			}
			stats.FilesSkipped++
			continue
		}

		// Check if already copied
		if _, exists := s.copied[file.Name()]; exists {
			s.log.Info("skipping duplicate", zap.String("file", file.Name()))
			stats.FilesSkipped++
			continue
		}

		// Copy the file
		srcPath := filepath.Join(playlistPath, file.Name())
		if err := s.copyFile(srcPath, file.Name()); err != nil {
			s.log.Error("failed to copy file",
				zap.String("file", file.Name()),
				zap.Error(err))
			stats.Errors = append(stats.Errors, fmt.Errorf("copy %s: %w", file.Name(), err))
			continue
		}

		s.copied[file.Name()] = struct{}{}
		stats.FilesCopied++
		s.log.Info("copied file", zap.String("file", file.Name()))
	}

	return nil
}

// copyFile copies a single file to the destination
func (s *Service) copyFile(srcPath, filename string) error {
	if s.cfg.DryRun {
		s.log.Info("dry-run: would copy file",
			zap.String("src", srcPath),
			zap.String("dst", filepath.Join(s.cfg.DestinationFolder, filename)))
		return nil
	}

	return utils.CopyFileToFolder(srcPath, s.cfg.DestinationFolder)
}

// RunApp is a compatibility wrapper for the old API
// Deprecated: Use NewService and Service.Run instead
func RunApp(log *zap.Logger, cfg *config.Config) (int, error) {
	svc := NewService(log, cfg)
	stats, err := svc.Run()
	if err != nil {
		return 0, err
	}
	return stats.FilesCopied, nil
}
