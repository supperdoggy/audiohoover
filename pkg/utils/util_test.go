package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileToFolder(t *testing.T) {
	// Create temp directories for testing
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("failed to create dst dir: %v", err)
	}

	tests := []struct {
		name         string
		setupFunc    func() (string, string)
		wantErr      bool
		validateFunc func(t *testing.T, srcFile, dstFolder string)
	}{
		{
			name: "successful copy",
			setupFunc: func() (string, string) {
				srcFile := filepath.Join(srcDir, "test.txt")
				if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				return srcFile, dstDir
			},
			wantErr: false,
			validateFunc: func(t *testing.T, srcFile, dstFolder string) {
				dstFile := filepath.Join(dstFolder, filepath.Base(srcFile))
				content, err := os.ReadFile(dstFile)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("content mismatch: got %q, want %q", string(content), "test content")
				}
			},
		},
		{
			name: "source file does not exist",
			setupFunc: func() (string, string) {
				return filepath.Join(srcDir, "nonexistent.txt"), dstDir
			},
			wantErr: true,
		},
		{
			name: "destination folder does not exist",
			setupFunc: func() (string, string) {
				srcFile := filepath.Join(srcDir, "test2.txt")
				if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				return srcFile, filepath.Join(tmpDir, "nonexistent")
			},
			wantErr: true,
		},
		{
			name: "copy file with spaces in name",
			setupFunc: func() (string, string) {
				srcFile := filepath.Join(srcDir, "test file with spaces.txt")
				if err := os.WriteFile(srcFile, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				return srcFile, dstDir
			},
			wantErr: false,
			validateFunc: func(t *testing.T, srcFile, dstFolder string) {
				dstFile := filepath.Join(dstFolder, filepath.Base(srcFile))
				if _, err := os.Stat(dstFile); err != nil {
					t.Errorf("destination file not found: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcFile, dstFolder := tt.setupFunc()
			err := CopyFileToFolder(srcFile, dstFolder)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFileToFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validateFunc != nil {
				tt.validateFunc(t, srcFile, dstFolder)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		setupFunc func() string
		want      bool
		wantErr   bool
	}{
		{
			name: "file exists",
			setupFunc: func() string {
				file := filepath.Join(tmpDir, "exists.txt")
				if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				return file
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file does not exist",
			setupFunc: func() string {
				return filepath.Join(tmpDir, "notexists.txt")
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "directory exists",
			setupFunc: func() string {
				dir := filepath.Join(tmpDir, "testdir")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("failed to create test dir: %v", err)
				}
				return dir
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()
			got, err := FileExists(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
