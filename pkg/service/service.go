package service

import (
	"os"
	"sync"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"github.com/supperdoggy/audiohoover/pkg/utils"
	"go.uber.org/zap"
)

type duplicatesMap struct {
	m   map[string]struct{}
	mut sync.Mutex
}

var (
	duplicates = duplicatesMap{
		m:   make(map[string]struct{}),
		mut: sync.Mutex{},
	}
)

func RunApp(log *zap.Logger, cfg *config.Config) (int, error) {
	// read destination folder and fill duplicates
	destinationFolder, err := os.Open(cfg.DestinationFolder)
	if err != nil {
		log.Fatal("failed to open destination folder", zap.Error(err))
	}

	existingFiles := getExistingFiles(log, destinationFolder)

	duplicates.m = existingFiles

	// read root folder and copy files
	playlistsFolder, err := os.Open(cfg.PlaylistsRoot)
	if err != nil {
		log.Fatal("failed to open root folder", zap.Error(err))
	}

	playlists, err := playlistsFolder.Readdir(0)
	if err != nil {
		log.Fatal("failed to read root folder", zap.Error(err))
	}

	for _, playlist := range playlists {
		path := cfg.PlaylistsRoot + "/" + playlist.Name()
		playlistFolder, err := os.Open(path)
		if err != nil {
			log.Fatal("failed to open playlist folder", zap.Error(err))
		}

		err = migratePlaylist(log, playlistFolder, path, cfg.DestinationFolder)
		if err != nil {
			log.Error("failed to migrate playlist", zap.Error(err))
			continue
		}
	}

	return len(duplicates.m), nil
}

// migratePlaylist copies files from playlist folder to destination folder
func migratePlaylist(log *zap.Logger, playlistFolder *os.File, playlistPath, destination string) error {
	files, err := playlistFolder.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		// check if not in duplicates
		if _, ok := duplicates.m[file.Name()]; ok {
			log.Info("found duplicate", zap.String("file", file.Name()))
			continue
		}

		// copy file to destination folder
		log.Info("copying file", zap.String("file", file.Name()))

		err := utils.CopyFileToFolder(playlistPath+"/"+file.Name(), destination)
		if err != nil {
			log.Error("failed to copy file", zap.Error(err))
			continue
		}

		duplicates.m[file.Name()] = struct{}{}
	}
	return nil
}

// getExistingFiles reads the destination folder and returns a map of existing files
func getExistingFiles(log *zap.Logger, destinationFolder *os.File) map[string]struct{} {
	duplicates := make(map[string]struct{})

	destinationFiles, err := destinationFolder.Readdir(0)
	if err != nil {
		log.Fatal("failed to read destination folder", zap.Error(err))
	}

	for _, file := range destinationFiles {
		duplicates[file.Name()] = struct{}{}
	}

	return duplicates
}
