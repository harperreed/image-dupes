package main

import (
	"reflect"
	"testing"

	"github.com/vitali-fedulov/images4"
)

func TestFindSimilarImages(t *testing.T) {
	tests := []struct {
		name     string
		images   []ImageInfo
		expected [][]string
	}{
		{
			name: "Identical images (same file hash)",
			images: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img2.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img3.jpg", FileHash: [16]byte{1, 2, 3, 4}},
			},
			expected: [][]string{{"img1.jpg", "img2.jpg", "img3.jpg"}},
		},
		{
			name: "Visually similar images (different file hash)",
			images: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}, Icon: images4.IconT{}},
				{Path: "img2.jpg", FileHash: [16]byte{2, 3, 4, 5}, Icon: images4.IconT{}},
				{Path: "img3.jpg", FileHash: [16]byte{3, 4, 5, 6}, Icon: images4.IconT{}},
			},
			expected: [][]string{{"img1.jpg", "img2.jpg", "img3.jpg"}},
		},
		{
			name: "Mixed set of identical, similar, and unique images",
			images: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img2.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img3.jpg", FileHash: [16]byte{2, 3, 4, 5}, Icon: images4.IconT{}},
				{Path: "img4.jpg", FileHash: [16]byte{3, 4, 5, 6}, Icon: images4.IconT{}},
				{Path: "img5.jpg", FileHash: [16]byte{4, 5, 6, 7}},
			},
			expected: [][]string{
				{"img1.jpg", "img2.jpg"},
				{"img3.jpg", "img4.jpg"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findSimilarImages(tt.images)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("findSimilarImages() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGroupByFileHash(t *testing.T) {
	tests := []struct {
		name     string
		images   []ImageInfo
		expected map[[16]byte][]ImageInfo
	}{
		{
			name: "Multiple images with same file hash",
			images: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img2.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img3.jpg", FileHash: [16]byte{2, 3, 4, 5}},
			},
			expected: map[[16]byte][]ImageInfo{
				{1, 2, 3, 4}: {
					{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}},
					{Path: "img2.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				},
				{2, 3, 4, 5}: {
					{Path: "img3.jpg", FileHash: [16]byte{2, 3, 4, 5}},
				},
			},
		},
		{
			name: "All images with unique file hashes",
			images: []ImageInfo{
				{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}},
				{Path: "img2.jpg", FileHash: [16]byte{2, 3, 4, 5}},
				{Path: "img3.jpg", FileHash: [16]byte{3, 4, 5, 6}},
			},
			expected: map[[16]byte][]ImageInfo{
				{1, 2, 3, 4}: {{Path: "img1.jpg", FileHash: [16]byte{1, 2, 3, 4}}},
				{2, 3, 4, 5}: {{Path: "img2.jpg", FileHash: [16]byte{2, 3, 4, 5}}},
				{3, 4, 5, 6}: {{Path: "img3.jpg", FileHash: [16]byte{3, 4, 5, 6}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := groupByFileHash(tt.images)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("groupByFileHash() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGroupByImageSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		images   []ImageInfo
		expected [][]string
	}{
		{
			name: "Visually similar images",
			images: []ImageInfo{
				{Path: "img1.jpg", Icon: images4.IconT{}},
				{Path: "img2.jpg", Icon: images4.IconT{}},
				{Path: "img3.jpg", Icon: images4.IconT{}},
			},
			expected: [][]string{{"img1.jpg", "img2.jpg", "img3.jpg"}},
		},
		{
			name: "Visually distinct images",
			images: []ImageInfo{
				{Path: "img1.jpg", Icon: images4.IconT{1}},
				{Path: "img2.jpg", Icon: images4.IconT{2}},
				{Path: "img3.jpg", Icon: images4.IconT{3}},
			},
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := groupByImageSimilarity(tt.images)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("groupByImageSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetRemainingImages(t *testing.T) {
	tests := []struct {
		name      string
		allImages []ImageInfo
		grouped   [][]string
		expected  []ImageInfo
	}{
		{
			name: "Some images already grouped",
			allImages: []ImageInfo{
				{Path: "img1.jpg"},
				{Path: "img2.jpg"},
				{Path: "img3.jpg"},
				{Path: "img4.jpg"},
			},
			grouped: [][]string{
				{"img1.jpg", "img2.jpg"},
			},
			expected: []ImageInfo{
				{Path: "img3.jpg"},
				{Path: "img4.jpg"},
			},
		},
		{
			name: "All images grouped",
			allImages: []ImageInfo{
				{Path: "img1.jpg"},
				{Path: "img2.jpg"},
			},
			grouped: [][]string{
				{"img1.jpg", "img2.jpg"},
			},
			expected: []ImageInfo{},
		},
		{
			name: "No images grouped",
			allImages: []ImageInfo{
				{Path: "img1.jpg"},
				{Path: "img2.jpg"},
			},
			grouped: [][]string{},
			expected: []ImageInfo{
				{Path: "img1.jpg"},
				{Path: "img2.jpg"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRemainingImages(tt.allImages, tt.grouped)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getRemainingImages() = %v, want %v", result, tt.expected)
			}
		})
	}
}
