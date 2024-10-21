package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestScanEmptyDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "empty_dir")
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
	tempDir, err := ioutil.TempDir("", "no_image_files")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create non-image files
	createTempFile(t, tempDir, "file1.txt")
	createTempFile(t, tempDir, "file2.pdf")

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(images))
	}
}

func TestScanOnlyImageFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "only_image_files")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create image files
	expectedImages := []string{
		createTempFile(t, tempDir, "image1.jpg"),
		createTempFile(t, tempDir, "image2.png"),
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
	tempDir, err := ioutil.TempDir("", "mixed_file_types")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mixed file types
	expectedImages := []string{
		createTempFile(t, tempDir, "image1.jpg"),
		createTempFile(t, tempDir, "image2.png"),
	}
	createTempFile(t, tempDir, "file1.txt")
	createTempFile(t, tempDir, "file2.pdf")

	images, err := scanDirectoryRecursive(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryRecursive failed: %v", err)
	}

	if !reflect.DeepEqual(sortStrings(images), sortStrings(expectedImages)) {
		t.Errorf("Expected images %v, got %v", expectedImages, images)
	}
}

func TestScanSubdirectories(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "subdirectories")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories with image files
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	os.Mkdir(subDir1, 0755)
	os.Mkdir(subDir2, 0755)

	expectedImages := []string{
		createTempFile(t, tempDir, "image1.jpg"),
		createTempFile(t, subDir1, "image2.png"),
		createTempFile(t, subDir2, "image3.jpeg"),
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
	tempDir, err := ioutil.TempDir("", "different_extensions")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create image files with different extensions
	expectedImages := []string{
		createTempFile(t, tempDir, "image1.jpg"),
		createTempFile(t, tempDir, "image2.jpeg"),
		createTempFile(t, tempDir, "image3.png"),
		createTempFile(t, tempDir, "image4.PNG"),
		createTempFile(t, tempDir, "image5.JPG"),
	}
	createTempFile(t, tempDir, "image6.gif") // This should not be included

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
	tempFile := createTempFile(t, "", "testfile.txt")
	defer os.Remove(tempFile)

	_, err = scanDirectoryRecursive(tempFile)
	if err == nil {
		t.Error("Expected an error when scanning a file instead of a directory, but got nil")
	}
}

// Helper function to create a temporary file and return its path
func createTempFile(t *testing.T, dir, name string) string {
	file, err := ioutil.TempFile(dir, name)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()
	return file.Name()
}

// Helper function to sort a slice of strings
func sortStrings(strs []string) []string {
	sorted := make([]string, len(strs))
	copy(sorted, strs)
	sort.Strings(sorted)
	return sorted
}
