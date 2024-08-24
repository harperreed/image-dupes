package main

import (
	// "image"
	_ "image/jpeg"
	_ "image/png"
		"fmt"
	// "os"
)

func findSimilarImages(imageInfos []ImageInfo) [][]string {
	var similarGroups [][]string
	totalComparisons := (len(imageInfos) * (len(imageInfos) - 1)) / 2
	comparisonsDone := 0

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

	// Pass 2: Pixel-by-pixel comparison (for images not grouped in Pass 1)
	remainingImages := getRemainingImages(imageInfos, similarGroups)
	pixelGroups := groupByPixelComparison(remainingImages, &comparisonsDone, totalComparisons)
	similarGroups = append(similarGroups, pixelGroups...)

	// Pass 3: Perceptual hash comparison (for images not grouped in Pass 1 or 2)
	remainingImages = getRemainingImages(imageInfos, similarGroups)
	pHashGroups := groupByPHash(remainingImages, &comparisonsDone, totalComparisons)
	similarGroups = append(similarGroups, pHashGroups...)

	return similarGroups
}

func groupByFileHash(imageInfos []ImageInfo) map[[16]byte][]ImageInfo {
	groups := make(map[[16]byte][]ImageInfo)
	for _, img := range imageInfos {
		groups[img.FileHash] = append(groups[img.FileHash], img)
	}
	return groups
}

func groupByPixelComparison(imageInfos []ImageInfo, comparisonsDone *int, totalComparisons int) [][]string {
	var groups [][]string
	compared := make(map[string]bool)

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

			if arePixelsSimilar(img1.Path, img2.Path) {
				group = append(group, img2.Path)
				compared[img2.Path] = true
			}

			*comparisonsDone++
			if *comparisonsDone%100 == 0 {
				fmt.Printf("\rComparing images: %d/%d", *comparisonsDone, totalComparisons)
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

func arePixelsSimilar(path1, path2 string) bool {
	img1, err1 := openImage(path1)
	img2, err2 := openImage(path2)
	if err1 != nil || err2 != nil {
		return false
	}

	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1 != bounds2 {
		return false
	}

	diffCount := 0
	threshold := bounds1.Dx() * bounds1.Dy() / 10 // Allow 10% difference

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			if abs(int(r1)-int(r2)) > 10000 || abs(int(g1)-int(g2)) > 10000 || abs(int(b1)-int(b2)) > 10000 {
				diffCount++
			}
			if diffCount > threshold {
				return false
			}
		}
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func groupByPHash(imageInfos []ImageInfo, comparisonsDone *int, totalComparisons int) [][]string {
	var groups [][]string
	compared := make(map[string]bool)

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

			distance, err := img1.PHash.Distance(img2.PHash)
			if err == nil && distance <= 10 {
				group = append(group, img2.Path)
				compared[img2.Path] = true
			}

			*comparisonsDone++
			if *comparisonsDone%100 == 0 {
				fmt.Printf("\rComparing images: %d/%d", *comparisonsDone, totalComparisons)
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
