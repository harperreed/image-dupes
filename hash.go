package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

type ImageHash struct {
	Path string
	Hash uint64
}

func computeHashes(images []string, progress *Progress) ([]ImageHash, error) {
	var imageHashes []ImageHash
	for _, path := range images {
		hash, err := computeDHash(path)
		if err != nil {
			fmt.Printf("\nError computing hash for %s: %v\n", path, err)
			continue
		}
		imageHashes = append(imageHashes, ImageHash{Path: path, Hash: hash})
		progress.Increment()
	}
	return imageHashes, nil
}

func computeDHash(path string) (uint64, error) {
	img, err := openImage(path)
	if err != nil {
		return 0, err
	}

	resized := resize.Resize(9, 8, img, resize.Bilinear)

	var hash uint64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			left, _ := resized.At(x, y).(color.Gray)
			right, _ := resized.At(x+1, y).(color.Gray)
			if left.Y > right.Y {
				hash |= 1 << (uint(y*8 + x))
			}
		}
	}

	return hash, nil
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}
