package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"github.com/supperdoggy/audiohoover/pkg/service"
	"go.uber.org/zap"
)

const version = "1.0.0"

var (
	playlistsRoot     = flag.String("playlists", "./playlists", "Root directory containing playlist folders")
	destinationFolder = flag.String("destination", "./output", "Destination folder for collected audio files")
	dryRun            = flag.Bool("dry-run", false, "Run without actually copying files")
	audioExtensions   = flag.String("extensions", ".mp3,.flac,.wav,.m4a,.aac,.ogg,.wma", "Comma-separated list of audio file extensions")
	verbose           = flag.Bool("verbose", false, "Enable verbose logging")
	showVersion       = flag.Bool("version", false, "Show version information")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("AudioHoover v%s\n", version)
		os.Exit(0)
	}

	// Initialize logger
	var log *zap.Logger
	var err error
	if *verbose {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Override environment with CLI flags
	if flag.Lookup("playlists").Value.String() != flag.Lookup("playlists").DefValue {
		os.Setenv("PLAYLISTS_ROOT", *playlistsRoot)
	}
	if flag.Lookup("destination").Value.String() != flag.Lookup("destination").DefValue {
		os.Setenv("DESTINATION_FOLDER", *destinationFolder)
	}
	if *dryRun {
		os.Setenv("DRY_RUN", "true")
	}
	if flag.Lookup("extensions").Value.String() != flag.Lookup("extensions").DefValue {
		os.Setenv("AUDIO_EXTENSIONS", *audioExtensions)
	}
	if *verbose {
		os.Setenv("VERBOSE", "true")
	}

	ctx := context.Background()
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	log.Info("starting AudioHoover",
		zap.String("version", version),
		zap.String("playlists_root", cfg.PlaylistsRoot),
		zap.String("destination", cfg.DestinationFolder),
		zap.Bool("dry_run", cfg.DryRun))

	svc := service.NewService(log, cfg)
	stats, err := svc.Run()
	if err != nil {
		log.Fatal("failed to run app", zap.Error(err))
	}

	// Report statistics
	log.Info("completed",
		zap.Int("files_processed", stats.FilesProcessed),
		zap.Int("files_copied", stats.FilesCopied),
		zap.Int("files_skipped", stats.FilesSkipped),
		zap.Int("errors", len(stats.Errors)))

	if len(stats.Errors) > 0 {
		log.Warn("encountered errors during processing")
		for i, err := range stats.Errors {
			log.Error("error", zap.Int("index", i+1), zap.Error(err))
		}
		os.Exit(1)
	}
}
