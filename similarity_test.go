package main

import (
	"testing"
)

func TestFindSimilarImages(t *testing.T) {
	imageInfos := []ImageInfo{
		{Path: "image1.jpg", FileHash: [16]byte{1}, Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image2.jpg", FileHash: [16]byte{1}, Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image3.jpg", FileHash: [16]byte{2}, Icon: images4.IconT{Pixels: []uint16{4, 5, 6}}},
		{Path: "image4.jpg", FileHash: [16]byte{3}, Icon: images4.IconT{Pixels: []uint16{7, 8, 9}}},
	}

	expectedGroups := [][]string{
		{"image1.jpg", "image2.jpg"},
		{"image3.jpg"},
		{"image4.jpg"},
	}

	similarGroups := findSimilarImages(imageInfos)

	if len(similarGroups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d", len(expectedGroups), len(similarGroups))
	}

	for i, group := range similarGroups {
		if len(group) != len(expectedGroups[i]) {
			t.Errorf("Expected group %d to have %d images, got %d", i, len(expectedGroups[i]), len(group))
		}
		for j, path := range group {
			if path != expectedGroups[i][j] {
				t.Errorf("Expected group %d image %d to be %s, got %s", i, j, expectedGroups[i][j], path)
			}
		}
	}
}

func TestGroupByFileHash(t *testing.T) {
	imageInfos := []ImageInfo{
		{Path: "image1.jpg", FileHash: [16]byte{1}},
		{Path: "image2.jpg", FileHash: [16]byte{1}},
		{Path: "image3.jpg", FileHash: [16]byte{2}},
	}

	expectedGroups := map[[16]byte][]ImageInfo{
		[16]byte{1}: {{Path: "image1.jpg", FileHash: [16]byte{1}}, {Path: "image2.jpg", FileHash: [16]byte{1}}},
		[16]byte{2}: {{Path: "image3.jpg", FileHash: [16]byte{2}}},
	}

	groups := groupByFileHash(imageInfos)

	if len(groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d", len(expectedGroups), len(groups))
	}

	for hash, group := range groups {
		if len(group) != len(expectedGroups[hash]) {
			t.Errorf("Expected group with hash %x to have %d images, got %d", hash, len(expectedGroups[hash]), len(group))
		}
		for i, img := range group {
			if img.Path != expectedGroups[hash][i].Path {
				t.Errorf("Expected group with hash %x image %d to be %s, got %s", hash, i, expectedGroups[hash][i].Path, img.Path)
			}
		}
	}
}

func TestGroupByImageSimilarity(t *testing.T) {
	imageInfos := []ImageInfo{
		{Path: "image1.jpg", Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image2.jpg", Icon: images4.IconT{Pixels: []uint16{1, 2, 3}}},
		{Path: "image3.jpg", Icon: images4.IconT{Pixels: []uint16{4, 5, 6}}},
	}

	expectedGroups := [][]string{
		{"image1.jpg", "image2.jpg"},
		{"image3.jpg"},
	}

	groups := groupByImageSimilarity(imageInfos)

	if len(groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d", len(expectedGroups), len(groups))
	}

	for i, group := range groups {
		if len(group) != len(expectedGroups[i]) {
			t.Errorf("Expected group %d to have %d images, got %d", i, len(expectedGroups[i]), len(group))
		}
		for j, path := range group {
			if path != expectedGroups[i][j] {
				t.Errorf("Expected group %d image %d to be %s, got %s", i, j, expectedGroups[i][j], path)
			}
		}
	}
}

func TestGetRemainingImages(t *testing.T) {
	allImages := []ImageInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
		{Path: "image3.jpg"},
	}

	groupedImages := [][]string{
		{"image1.jpg", "image2.jpg"},
	}

	expectedRemaining := []ImageInfo{
		{Path: "image3.jpg"},
	}

	remaining := getRemainingImages(allImages, groupedImages)

	if len(remaining) != len(expectedRemaining) {
		t.Errorf("Expected %d remaining images, got %d", len(expectedRemaining), len(remaining))
	}

	for i, img := range remaining {
		if img.Path != expectedRemaining[i].Path {
			t.Errorf("Expected remaining image %d to be %s, got %s", i, expectedRemaining[i].Path, img.Path)
		}
	}
}
