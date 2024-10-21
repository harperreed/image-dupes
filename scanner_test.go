package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestScanEmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "empty_dir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(images))
	}
}

func TestScanNoImageFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "no_image_files")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create non-image files
	createNamedTempFile(t, tempDir, "file1.txt")
	createNamedTempFile(t, tempDir, "file2.pdf")

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(images))
	}
}

func TestScanOnlyImageFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "only_image_files")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create image files
	expectedImages := []string{
		createNamedTempFile(t, tempDir, "image1.jpg"),
		createNamedTempFile(t, tempDir, "image2.png"),
	}

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if !reflect.DeepEqual(sortStrings(images), sortStrings(expectedImages)) {
		t.Errorf("Expected images %v, got %v", expectedImages, images)
	}
}

func TestScanMixedFileTypes(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mixed_file_types")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mixed file types
	expectedImages := []string{
		createNamedTempFile(t, tempDir, "image1.jpg"),
		createNamedTempFile(t, tempDir, "image2.png"),
	}
	createNamedTempFile(t, tempDir, "file1.txt")
	createNamedTempFile(t, tempDir, "file2.pdf")

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if !reflect.DeepEqual(sortStrings(images), sortStrings(expectedImages)) {
		t.Errorf("Expected images %v, got %v", expectedImages, images)
	}
}

func TestScanSubdirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subdirectories")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories with image files
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	if err := os.Mkdir(subDir1, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	if err := os.Mkdir(subDir2, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	expectedImages := []string{
		createNamedTempFile(t, tempDir, "image1.jpg"),
		createNamedTempFile(t, subDir1, "image2.png"),
		createNamedTempFile(t, subDir2, "image3.jpeg"),
	}

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if !reflect.DeepEqual(sortStrings(images), sortStrings(expectedImages)) {
		t.Errorf("Expected images %v, got %v", expectedImages, images)
	}
}

func TestScanDifferentImageExtensions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "different_extensions")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create image files with different extensions
	expectedImages := []string{
		createNamedTempFile(t, tempDir, "image1.jpg"),
		createNamedTempFile(t, tempDir, "image2.jpeg"),
		createNamedTempFile(t, tempDir, "image3.png"),
		createNamedTempFile(t, tempDir, "image4.PNG"),
		createNamedTempFile(t, tempDir, "image5.JPG"),
	}
	createNamedTempFile(t, tempDir, "image6.gif") // This should not be included

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if !reflect.DeepEqual(sortStrings(images), sortStrings(expectedImages)) {
		t.Errorf("Expected images %v, got %v", expectedImages, images)
	}
}

func TestScanErrorHandling(t *testing.T) {
	// Test with a non-existent directory
	_, err := scanDirectoryRecursive("/path/to/nonexistent/directory")
	if err == nil {
		t.Error("Expected an error for non-existent directory, but got nil")
	}

	// Test with a file instead of a directory
	tempFile := createNamedTempFile(t, "", "testfile.txt")
	defer os.Remove(tempFile)

	_, err = scanDirectoryRecursive(tempFile)
	if err == nil {
		t.Error("Expected an error when scanning a file instead of a directory, but got nil")
	}
}

// Helper function to create a temporary file and return its path
func createNamedTempFile(t *testing.T, dir, name string) string {
	filePath := filepath.Join(dir, name)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()
	return filePath
}

// Helper function to sort a slice of strings
func sortStrings(strs []string) []string {
	sorted := make([]string, len(strs))
	copy(sorted, strs)
	sort.Strings(sorted)
	return sorted
}
