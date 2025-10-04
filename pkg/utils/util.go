package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFileToFolder copies a file from srcFile to dstFolder, preserving the filename.
// It returns an error if the operation fails at any step.
func CopyFileToFolder(srcFile, dstFolder string) error {
	// Extract the file name from the source path
	fileName := filepath.Base(srcFile)

	// Create the destination path
	dstFile := filepath.Join(dstFolder, fileName)

	// Open the source file
	source, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	// Create the destination file
	destination, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	// Copy the contents from source to destination
	if _, err = io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Ensure the file is flushed to disk
	if err = destination.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// FileExists checks if a file exists at the given path
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
