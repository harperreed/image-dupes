package main

import (
	"fmt"

	"github.com/vitali-fedulov/images4"
)

func findSimilarImages(imageInfos []ImageInfo) [][]string {
	var similarGroups [][]string

	// Pass 1: File hash comparison
	fileHashGroups := groupByFileHash(imageInfos)
	for _, group := range fileHashGroups {
		if len(group) > 1 {
			var paths []string
			for _, img := range group {
				paths = append(paths, img.Path)
			}
			similarGroups = append(similarGroups, paths)
		}
	}

	// Pass 2: Image comparison
	remainingImages := getRemainingImages(imageInfos, similarGroups)
	imgGroups := groupByImageSimilarity(remainingImages)
	similarGroups = append(similarGroups, imgGroups...)

	return similarGroups
}

func groupByFileHash(imageInfos []ImageInfo) map[[16]byte][]ImageInfo {
	groups := make(map[[16]byte][]ImageInfo)
	for _, img := range imageInfos {
		groups[img.FileHash] = append(groups[img.FileHash], img)
	}
	return groups
}

func groupByImageSimilarity(imageInfos []ImageInfo) [][]string {
	var groups [][]string
	compared := make(map[string]bool)
	totalComparisons := (len(imageInfos) * (len(imageInfos) - 1)) / 2
	comparisonsDone := 0

	for i, img1 := range imageInfos {
		if compared[img1.Path] {
			continue
		}

		group := []string{img1.Path}
		for j := i + 1; j < len(imageInfos); j++ {
			img2 := imageInfos[j]
			if compared[img2.Path] {
				continue
			}

			if images4.Similar(img1.Icon, img2.Icon) {
				group = append(group, img2.Path)
				compared[img2.Path] = true
			}

			comparisonsDone++
			if comparisonsDone%10 == 0 || comparisonsDone == totalComparisons {
				fmt.Printf("\rImage comparison progress: %d/%d", comparisonsDone, totalComparisons)
			}
		}

		if len(group) > 1 {
			groups = append(groups, group)
		}
		compared[img1.Path] = true
	}

	fmt.Println() // New line after progress
	return groups
}

func getRemainingImages(allImages []ImageInfo, groupedImages [][]string) []ImageInfo {
	grouped := make(map[string]bool)
	for _, group := range groupedImages {
		for _, path := range group {
			grouped[path] = true
		}
	}

	var remaining []ImageInfo
	for _, img := range allImages {
		if !grouped[img.Path] {
			remaining = append(remaining, img)
		}
	}
	return remaining
}
