package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"image"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/vitali-fedulov/images4"
)

// Mock structs
type MockImageOpener struct{}

func (m MockImageOpener) Open(path string) (image.Image, error) {
	if path == "nonexistent.jpg" {
		return nil, errors.New("file not found")
	}
	return &image.RGBA{}, nil
}

type MockIconCreator struct{}

func (m MockIconCreator) Icon(img image.Image) images4.IconT {
	return images4.IconT{Pixels: []uint16{1, 2, 3}} // Non-empty IconT
}

type MockFileHasher struct{}

func (m MockFileHasher) ComputeFileHash(path string) ([16]byte, error) {
	if path == "nonexistent.jpg" {
		return [16]byte{}, errors.New("file not found")
	}
	return md5.Sum([]byte(path)), nil // Use path as content for deterministic testing
}

// Helper function to create a temporary file with content
func createTempFile(t *testing.T, content []byte) string {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	return tmpfile.Name()
}

// Helper function to remove temporary files
func removeTempFiles(t *testing.T, files []string) {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			t.Errorf("Failed to remove temp file %s: %v", f, err)
		}
	}
}

func TestComputeHashes(t *testing.T) {
	// Create temporary test files
	content1 := []byte("test content 1")
	content2 := []byte("test content 2")
	file1 := createTempFile(t, content1)
	file2 := createTempFile(t, content2)
	defer removeTempFiles(t, []string{file1, file2})

	// Test cases
	testCases := []struct {
		name                    string
		imagePaths              []string
		expectedCount           int
		expectedProgressUpdates int
	}{
		{"SingleValidFile", []string{file1}, 1, 1},
		{"MultipleValidFiles", []string{file1, file2}, 2, 2},
		{"MixedValidAndInvalid", []string{file1, "nonexistent.jpg", file2}, 2, 3},
		{"EmptyList", []string{}, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			progress := make(chan string, len(tc.imagePaths))

			imageInfos, err := computeHashes(tc.imagePaths, progress, MockImageOpener{}, MockIconCreator{}, MockFileHasher{})

			if err != nil {
				t.Errorf("computeHashes returned an error: %v", err)
				return
			}

			if len(imageInfos) != tc.expectedCount {
				t.Errorf("Expected %d ImageInfo structs, got %d", tc.expectedCount, len(imageInfos))
			}

			// Check if FileHash and Icon fields are populated
			for _, info := range imageInfos {
				if info.FileHash == [16]byte{} {
					t.Errorf("FileHash is empty for %s", info.Path)
				}
				// Check if Icon is the zero value of images4.IconT
				if reflect.DeepEqual(info.Icon, images4.IconT{}) {
					t.Errorf("Icon is empty for %s", info.Path)
				}
			}

			// Check progress channel
			progressCount := 0
			if tc.expectedProgressUpdates > 0 {
				timeout := time.After(1 * time.Second)
			progressLoop:
				for {
					select {
					case _, ok := <-progress:
						if !ok {
							break progressLoop
						}
						progressCount++
						if progressCount == tc.expectedProgressUpdates {
							break progressLoop
						}
					case <-timeout:
						t.Errorf("Test timed out waiting for progress updates")
						break progressLoop
					}
				}
			}

			if progressCount != tc.expectedProgressUpdates {
				t.Errorf("Expected %d progress updates, got %d", tc.expectedProgressUpdates, progressCount)
			}

			// Ensure the progress channel is empty
			select {
			case update, ok := <-progress:
				if ok {
					t.Errorf("Unexpected progress update: %s", update)
				}
			default:
				// Channel is empty, which is what we want
			}
		})
	}
}

func TestDefaultFileHasher(t *testing.T) {
	hasher := DefaultFileHasher{}

	// Test cases
	testCases := []struct {
		name        string
		content     []byte
		expectError bool
	}{
		{"ValidFile", []byte("test content"), false},
		{"EmptyFile", []byte{}, false},
		{"LargeFile", bytes.Repeat([]byte("a"), 1024*1024), false}, // 1MB file
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpfile := createTempFile(t, tc.content)
			defer removeTempFiles(t, []string{tmpfile})

			hash, err := hasher.ComputeFileHash(tmpfile)

			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !tc.expectError {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				expectedHash := md5.Sum(tc.content)
				if hash != expectedHash {
					t.Errorf("Expected hash %x, got %x", expectedHash, hash)
				}
			}
		})
	}

	// Test with non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := hasher.ComputeFileHash("nonexistent.file")
		if err == nil {
			t.Errorf("Expected an error for non-existent file, but got none")
		}
	})

	// Test with a directory instead of a file
	t.Run("Directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		_, err = hasher.ComputeFileHash(tmpDir)
		if err == nil {
			t.Errorf("Expected an error for directory, but got none")
		}
	})
}
