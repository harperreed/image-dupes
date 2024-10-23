package main

import (
	"fmt"
	"testing"

	"github.com/vitali-fedulov/images4"
)

// Mock for images4.Similar function
var mockSimilar func(icon1, icon2 images4.IconT) bool

func init() {
	// Replace the original Similar function with our mock
	images4.Similar = func(icon1, icon2 images4.IconT) bool {
		return mockSimilar(icon1, icon2)
	}
}

// Helper function to set up the mock behavior
func setMockSimilar(similar bool) {
	mockSimilar = func(icon1, icon2 images4.IconT) bool {
		return similar
	}
}

// ImageInfo is a struct that holds information about an image
type ImageInfo struct {
	Path     string
	FileHash [16]byte
	Icon     images4.IconT
}

func TestFindSimilarImages(t *testing.T) {
	// Test with a set of identical images (same file hash)
	identicalImages := []ImageInfo{
		{Path: "image1.jpg", FileHash: [16]byte{1}},
		{Path: "image2.jpg", FileHash: [16]byte{1}},
	}
	similarGroups := findSimilarImages(identicalImages)
	if len(similarGroups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(similarGroups))
	}

	// Test with a set of visually similar images (different file hash)
	similarImages := []ImageInfo{
		{Path: "image3.jpg", FileHash: [16]byte{2}, Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image4.jpg", FileHash: [16]byte{3}, Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
	}
	similarGroups = findSimilarImages(similarImages)
	if len(similarGroups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(similarGroups))
	}

	// Test with a mixed set of identical, similar, and unique images
	mixedImages := []ImageInfo{
		{Path: "image5.jpg", FileHash: [16]byte{4}},
		{Path: "image6.jpg", FileHash: [16]byte{4}},
		{Path: "image7.jpg", FileHash: [16]byte{5}, Icon: images4.IconT{Pixels: []uint16{4, 5, 6}}},
		{Path: "image8.jpg", FileHash: [16]byte{6}, Icon: images4.IconT{Pixels: []uint16{4, 5, 6}}},
		{Path: "image9.jpg", FileHash: [16]byte{7}},
	}
	similarGroups = findSimilarImages(mixedImages)
	if len(similarGroups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(similarGroups))
	}
}

func TestGroupByFileHash(t *testing.T) {
	// Test with multiple images having the same file hash
	images := []ImageInfo{
		{Path: "image1.jpg", FileHash: [16]byte{1}},
		{Path: "image2.jpg", FileHash: [16]byte{1}},
		{Path: "image3.jpg", FileHash: [16]byte{2}},
	}
	groups := groupByFileHash(images)
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	// Test with all images having unique file hashes
	uniqueImages := []ImageInfo{
		{Path: "image4.jpg", FileHash: [16]byte{3}},
		{Path: "image5.jpg", FileHash: [16]byte{4}},
	}
	groups = groupByFileHash(uniqueImages)
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

func TestGroupByImageSimilarity(t *testing.T) {
	// Test with visually similar images
	setMockSimilar(true)
	similarImages := []ImageInfo{
		{Path: "image1.jpg", Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image2.jpg", Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
	}
	groups := groupByImageSimilarity(similarImages)
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}

	// Test with visually distinct images
	setMockSimilar(false)
	distinctImages := []ImageInfo{
		{Path: "image3.jpg", Icon: images4.IconT{Pixels: []uint16{4, 5, 6}}},
		{Path: "image4.jpg", Icon: images4.IconT{Pixels: []uint16{7, 8, 9}}},
	}
	groups = groupByImageSimilarity(distinctImages)
	if len(groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(groups))
	}

	// Test performance with a large number of images
	setMockSimilar(true)
	var largeImages []ImageInfo
	for i := 0; i < 1000; i++ {
		largeImages = append(largeImages, ImageInfo{Path: fmt.Sprintf("image%d.jpg", i), Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}})
	}
	groups = groupByImageSimilarity(largeImages)
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}
}

func TestGetRemainingImages(t *testing.T) {
	// Test with all images grouped
	allImages := []ImageInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
	}
	groupedImages := [][]string{
		{"image1.jpg", "image2.jpg"},
	}
	remainingImages := getRemainingImages(allImages, groupedImages)
	if len(remainingImages) != 0 {
		t.Errorf("Expected 0 remaining images, got %d", len(remainingImages))
	}

	// Test with some images grouped and some remaining
	someImages := []ImageInfo{
		{Path: "image3.jpg"},
		{Path: "image4.jpg"},
		{Path: "image5.jpg"},
	}
	groupedImages = [][]string{
		{"image3.jpg", "image4.jpg"},
	}
	remainingImages = getRemainingImages(someImages, groupedImages)
	if len(remainingImages) != 1 {
		t.Errorf("Expected 1 remaining image, got %d", len(remainingImages))
	}

	// Test with no images grouped
	noGroupedImages := []ImageInfo{
		{Path: "image6.jpg"},
		{Path: "image7.jpg"},
	}
	groupedImages = [][]string{}
	remainingImages = getRemainingImages(noGroupedImages, groupedImages)
	if len(remainingImages) != 2 {
		t.Errorf("Expected 2 remaining images, got %d", len(remainingImages))
	}
}
