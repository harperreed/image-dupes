package main

import (
	"math/bits"
)

func findSimilarImages(imageHashes []ImageHash, threshold int) [][]string {
	var similarGroups [][]string

	for i := 0; i < len(imageHashes); i++ {
		group := []string{imageHashes[i].Path}

		for j := i + 1; j < len(imageHashes); j++ {
			distance := hammingDistance(imageHashes[i].Hash, imageHashes[j].Hash)

			if distance <= threshold {
				group = append(group, imageHashes[j].Path)
				imageHashes = append(imageHashes[:j], imageHashes[j+1:]...)
				j--
			}
		}

		if len(group) > 1 {
			similarGroups = append(similarGroups, group)
		}
	}

	return similarGroups
}

func hammingDistance(hash1, hash2 uint64) int {
	return bits.OnesCount64(hash1 ^ hash2)
}
