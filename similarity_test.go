package main

import (
	"reflect"
	"testing"

	"github.com/vitali-fedulov/images4"
)

// MockImageInfo is a mock struct for ImageInfo
type MockImageInfo struct {
	Path     string
	FileHash [16]byte
	Icon     images4.IconT
}

// mockSimilar is a mock function for images4.Similar
func mockSimilar(icon1, icon2 images4.IconT) bool {
	// For testing purposes, we'll consider icons similar if their first pixel is the same
	return icon1.Pixels[0] == icon2.Pixels[0]
}

// createMockImageInfos creates a slice of MockImageInfo for testing
func createMockImageInfos(paths []string, hashes [][16]byte, icons []images4.IconT) []MockImageInfo {
	infos := make([]MockImageInfo, len(paths))
	for i := range paths {
		infos[i] = MockImageInfo{
			Path:     paths[i],
			FileHash: hashes[i],
			Icon:     icons[i],
		}
	}
	return infos
}

// Helper function to create a mock IconT
func createMockIcon(pixel uint16) images4.IconT {
	return images4.IconT{Pixels: []uint16{pixel}}
}

func TestFindSimilarImages(t *testing.T) {
	// Override the images4.Similar function with our mock
	originalSimilar := images4.Similar
	images4.Similar = mockSimilar
	defer func() { images4.Similar = originalSimilar }()

	tests := []struct {
		name     string
		infos    []MockImageInfo
		expected [][]string
	}{
		{
			name: "Identical images (same file hash)",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {1}, {1}},
				[]images4.IconT{createMockIcon(1), createMockIcon(1), createMockIcon(1)},
			),
			expected: [][]string{{"img1.jpg", "img2.jpg", "img3.jpg"}},
		},
		{
			name: "Visually similar images (different file hash)",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(1), createMockIcon(2)},
			),
			expected: [][]string{{"img1.jpg", "img2.jpg"}},
		},
		{
			name: "Mixed set of identical, similar, and unique images",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg", "img4.jpg", "img5.jpg"},
				[][16]byte{{1}, {1}, {2}, {3}, {4}},
				[]images4.IconT{createMockIcon(1), createMockIcon(1), createMockIcon(2), createMockIcon(2), createMockIcon(3)},
			),
			expected: [][]string{{"img1.jpg", "img2.jpg"}, {"img3.jpg", "img4.jpg"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert MockImageInfo to ImageInfo
			imageInfos := make([]ImageInfo, len(tt.infos))
			for i, info := range tt.infos {
				imageInfos[i] = ImageInfo(info)
			}

			result := findSimilarImages(imageInfos)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("findSimilarImages() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGroupByFileHash(t *testing.T) {
	tests := []struct {
		name     string
		infos    []MockImageInfo
		expected map[[16]byte][]ImageInfo
	}{
		{
			name: "Multiple images with same file hash",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {1}, {2}},
				[]images4.IconT{createMockIcon(1), createMockIcon(1), createMockIcon(2)},
			),
			expected: map[[16]byte][]ImageInfo{
				{1}: {
					{Path: "img1.jpg", FileHash: [16]byte{1}, Icon: createMockIcon(1)},
					{Path: "img2.jpg", FileHash: [16]byte{1}, Icon: createMockIcon(1)},
				},
				{2}: {
					{Path: "img3.jpg", FileHash: [16]byte{2}, Icon: createMockIcon(2)},
				},
			},
		},
		{
			name: "All images with unique file hashes",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(2), createMockIcon(3)},
			),
			expected: map[[16]byte][]ImageInfo{
				{1}: {{Path: "img1.jpg", FileHash: [16]byte{1}, Icon: createMockIcon(1)}},
				{2}: {{Path: "img2.jpg", FileHash: [16]byte{2}, Icon: createMockIcon(2)}},
				{3}: {{Path: "img3.jpg", FileHash: [16]byte{3}, Icon: createMockIcon(3)}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert MockImageInfo to ImageInfo
			imageInfos := make([]ImageInfo, len(tt.infos))
			for i, info := range tt.infos {
				imageInfos[i] = ImageInfo(info)
			}

			result := groupByFileHash(imageInfos)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("groupByFileHash() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGroupByImageSimilarity(t *testing.T) {
	// Override the images4.Similar function with our mock
	originalSimilar := images4.Similar
	images4.Similar = mockSimilar
	defer func() { images4.Similar = originalSimilar }()

	tests := []struct {
		name     string
		infos    []MockImageInfo
		expected [][]string
	}{
		{
			name: "Visually similar images",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(1), createMockIcon(2)},
			),
			expected: [][]string{{"img1.jpg", "img2.jpg"}},
		},
		{
			name: "Visually distinct images",
			infos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(2), createMockIcon(3)},
			),
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert MockImageInfo to ImageInfo
			imageInfos := make([]ImageInfo, len(tt.infos))
			for i, info := range tt.infos {
				imageInfos[i] = ImageInfo(info)
			}

			result := groupByImageSimilarity(imageInfos)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("groupByImageSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkGroupByImageSimilarity(b *testing.B) {
	// Create a large set of mock image infos
	const numImages = 1000
	infos := make([]MockImageInfo, numImages)
	for i := 0; i < numImages; i++ {
		infos[i] = MockImageInfo{
			Path:     fmt.Sprintf("img%d.jpg", i),
			FileHash: [16]byte{byte(i % 256)},
			Icon:     createMockIcon(uint16(i % 5)), // This will create groups of similar images
		}
	}

	// Convert MockImageInfo to ImageInfo
	imageInfos := make([]ImageInfo, len(infos))
	for i, info := range infos {
		imageInfos[i] = ImageInfo(info)
	}

	// Override the images4.Similar function with our mock
	originalSimilar := images4.Similar
	images4.Similar = mockSimilar
	defer func() { images4.Similar = originalSimilar }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		groupByImageSimilarity(imageInfos)
	}
}

func TestGetRemainingImages(t *testing.T) {
	tests := []struct {
		name     string
		allInfos []MockImageInfo
		grouped  [][]string
		expected []ImageInfo
	}{
		{
			name: "All images grouped",
			allInfos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(2), createMockIcon(3)},
			),
			grouped:  [][]string{{"img1.jpg", "img2.jpg", "img3.jpg"}},
			expected: []ImageInfo{},
		},
		{
			name: "Some images grouped, some remaining",
			allInfos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg", "img4.jpg"},
				[][16]byte{{1}, {2}, {3}, {4}},
				[]images4.IconT{createMockIcon(1), createMockIcon(2), createMockIcon(3), createMockIcon(4)},
			),
			grouped: [][]string{{"img1.jpg", "img2.jpg"}},
			expected: []ImageInfo{
				{Path: "img3.jpg", FileHash: [16]byte{3}, Icon: createMockIcon(3)},
				{Path: "img4.jpg", FileHash: [16]byte{4}, Icon: createMockIcon(4)},
			},
		},
		{
			name: "No images grouped",
			allInfos: createMockImageInfos(
				[]string{"img1.jpg", "img2.jpg", "img3.jpg"},
				[][16]byte{{1}, {2}, {3}},
				[]images4.IconT{createMockIcon(1), createMockIcon(2), createMockIcon(3)},
			),
			grouped: [][]string{},
			expected: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1}, Icon: createMockIcon(1)},
				{Path: "img2.jpg", FileHash: [16]byte{2}, Icon: createMockIcon(2)},
				{Path: "img3.jpg", FileHash: [16]byte{3}, Icon: createMockIcon(3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert MockImageInfo to ImageInfo
			allImageInfos := make([]ImageInfo, len(tt.allInfos))
			for i, info := range tt.allInfos {
				allImageInfos[i] = ImageInfo(info)
			}

			result := getRemainingImages(allImageInfos, tt.grouped)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getRemainingImages() = %v, want %v", result, tt.expected)
			}
		})
	}
}
