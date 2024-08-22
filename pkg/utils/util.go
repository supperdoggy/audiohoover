package utils

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFileToFolder(srcFile, dstFolder string) error {
	// Extract the file name from the source path
	fileName := filepath.Base(srcFile)

	// Create the destination path
	dstFile := filepath.Join(dstFolder, fileName)

	// Open the source file
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create the destination file
	destination, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Ensure the file is flushed to disk
	err = destination.Sync()
	if err != nil {
		return err
	}

	return nil
}
