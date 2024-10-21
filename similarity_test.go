package main

import (
	"reflect"
	"testing"

	"github.com/vitali-fedulov/images4"
)

// MockSimilarFunc is a type for mocking the Similar function
type MockSimilarFunc func(icon1, icon2 images4.IconT) bool

// mockSimilar is a variable to hold the mock function
var mockSimilar MockSimilarFunc

// similarWrapper wraps the Similar function to allow mocking in tests
func similarWrapper(icon1, icon2 images4.IconT) bool {
	if mockSimilar != nil {
		return mockSimilar(icon1, icon2)
	}
	return images4.Similar(icon1, icon2)
}

func TestFindSimilarImages(t *testing.T) {
	// Create test data
	img1 := ImageInfo{Path: "img1.jpg", FileHash: [16]byte{1}, Icon: images4.IconT{}}
	img2 := ImageInfo{Path: "img2.jpg", FileHash: [16]byte{1}, Icon: images4.IconT{}}
	img3 := ImageInfo{Path: "img3.jpg", FileHash: [16]byte{2}, Icon: images4.IconT{}}
	img4 := ImageInfo{Path: "img4.jpg", FileHash: [16]byte{3}, Icon: images4.IconT{}}

	// Test case 1: Identical images (same file hash)
	t.Run("Identical Images", func(t *testing.T) {
		imageInfos := []ImageInfo{img1, img2}
		result := findSimilarImages(imageInfos)
		expected := [][]string{{"img1.jpg", "img2.jpg"}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test case 2: Visually similar images (different file hash)
	t.Run("Visually Similar Images", func(t *testing.T) {
		imageInfos := []ImageInfo{img3, img4}
		// Mock the Similar function to return true for these images
		mockSimilar = func(icon1, icon2 images4.IconT) bool {
			return true
		}
		defer func() { mockSimilar = nil }()

		result := findSimilarImages(imageInfos)
		expected := [][]string{{"img3.jpg", "img4.jpg"}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test case 3: Mixed set of identical, similar, and unique images
	t.Run("Mixed Image Set", func(t *testing.T) {
		imageInfos := []ImageInfo{img1, img2, img3, img4}
		// Mock the Similar function to return true only for img3 and img4
		mockSimilar = func(icon1, icon2 images4.IconT) bool {
			return (icon1 == img3.Icon && icon2 == img4.Icon) || (icon1 == img4.Icon && icon2 == img3.Icon)
		}
		defer func() { mockSimilar = nil }()

		result := findSimilarImages(imageInfos)
		expected := [][]string{{"img1.jpg", "img2.jpg"}, {"img3.jpg", "img4.jpg"}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestGroupByFileHash(t *testing.T) {
	// Test case 1: Multiple images with the same file hash
	t.Run("Same File Hash", func(t *testing.T) {
		imageInfos := []ImageInfo{
			{Path: "img1.jpg", FileHash: [16]byte{1}},
			{Path: "img2.jpg", FileHash: [16]byte{1}},
			{Path: "img3.jpg", FileHash: [16]byte{1}},
		}
		result := groupByFileHash(imageInfos)
		if len(result) != 1 {
			t.Errorf("Expected 1 group, got %d", len(result))
		}
		if len(result[[16]byte{1}]) != 3 {
			t.Errorf("Expected 3 images in group, got %d", len(result[[16]byte{1}]))
		}
	})

	// Test case 2: All images with unique file hashes
	t.Run("Unique File Hashes", func(t *testing.T) {
		imageInfos := []ImageInfo{
			{Path: "img1.jpg", FileHash: [16]byte{1}},
			{Path: "img2.jpg", FileHash: [16]byte{2}},
			{Path: "img3.jpg", FileHash: [16]byte{3}},
		}
		result := groupByFileHash(imageInfos)
		if len(result) != 3 {
			t.Errorf("Expected 3 groups, got %d", len(result))
		}
		for _, group := range result {
			if len(group) != 1 {
				t.Errorf("Expected 1 image in each group, got %d", len(group))
			}
		}
	})
}

func TestGroupByImageSimilarity(t *testing.T) {
	// Test case 1: Visually similar images
	t.Run("Similar Images", func(t *testing.T) {
		imageInfos := []ImageInfo{
			{Path: "img1.jpg", Icon: images4.IconT{}},
			{Path: "img2.jpg", Icon: images4.IconT{}},
			{Path: "img3.jpg", Icon: images4.IconT{}},
		}
		// Mock the Similar function to return true for all comparisons
		mockSimilar = func(icon1, icon2 images4.IconT) bool {
			return true
		}
		defer func() { mockSimilar = nil }()

		result := groupByImageSimilarity(imageInfos)
		if len(result) != 1 {
			t.Errorf("Expected 1 group, got %d", len(result))
		}
		if len(result[0]) != 3 {
			t.Errorf("Expected 3 images in group, got %d", len(result[0]))
		}
	})

	// Test case 2: Visually distinct images
	t.Run("Distinct Images", func(t *testing.T) {
		imageInfos := []ImageInfo{
			{Path: "img1.jpg", Icon: images4.IconT{}},
			{Path: "img2.jpg", Icon: images4.IconT{}},
			{Path: "img3.jpg", Icon: images4.IconT{}},
		}
		// Mock the Similar function to return false for all comparisons
		mockSimilar = func(icon1, icon2 images4.IconT) bool {
			return false
		}
		defer func() { mockSimilar = nil }()

		result := groupByImageSimilarity(imageInfos)
		if len(result) != 0 {
			t.Errorf("Expected 0 groups, got %d", len(result))
		}
	})
}

func TestGetRemainingImages(t *testing.T) {
	// Test case 1: Some images already grouped
	t.Run("Partially Grouped Images", func(t *testing.T) {
		allImages := []ImageInfo{
			{Path: "img1.jpg"},
			{Path: "img2.jpg"},
			{Path: "img3.jpg"},
			{Path: "img4.jpg"},
		}
		groupedImages := [][]string{
			{"img1.jpg", "img2.jpg"},
		}
		result := getRemainingImages(allImages, groupedImages)
		if len(result) != 2 {

			t.Errorf("Expected 2 remaining images, got %d", len(result))
		}
		expectedPaths := []string{"img3.jpg", "img4.jpg"}
		for i, img := range result {
			if img.Path != expectedPaths[i] {
				t.Errorf("Expected %s, got %s", expectedPaths[i], img.Path)
			}
		}
	})

	// Test case 2: All images grouped
	t.Run("All Images Grouped", func(t *testing.T) {
		allImages := []ImageInfo{
			{Path: "img1.jpg"},
			{Path: "img2.jpg"},
		}
		groupedImages := [][]string{
			{"img1.jpg", "img2.jpg"},
		}
		result := getRemainingImages(allImages, groupedImages)
		if len(result) != 0 {
			t.Errorf("Expected 0 remaining images, got %d", len(result))
		}
	})

	// Test case 3: No images grouped
	t.Run("No Images Grouped", func(t *testing.T) {
		allImages := []ImageInfo{
			{Path: "img1.jpg"},
			{Path: "img2.jpg"},
		}
		var groupedImages [][]string
		result := getRemainingImages(allImages, groupedImages)
		if len(result) != 2 {
			t.Errorf("Expected 2 remaining images, got %d", len(result))
		}
	})
}
