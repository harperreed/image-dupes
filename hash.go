package main

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
)

type ImageInfo struct {
	Path     string
	FileHash [16]byte
	PHash    *goimagehash.ImageHash
}

func computeHashes(images []string, progress chan<- string) ([]ImageInfo, error) {
	var imageInfos []ImageInfo
	for i, path := range images {
		fileHash, err := computeFileHash(path)
		if err != nil {
			fmt.Printf("Error computing file hash for %s: %v\n", path, err)
			continue
		}

		pHash, err := computePHash(path)
		if err != nil {
			fmt.Printf("Error computing perceptual hash for %s: %v\n", path, err)
			continue
		}

		imageInfos = append(imageInfos, ImageInfo{Path: path, FileHash: fileHash, PHash: pHash})

		// Send progress update
		progress <- fmt.Sprintf("Processed %d/%d: %s", i+1, len(images), filepath.Base(path))
	}
	return imageInfos, nil
}

func computeFileHash(path string) ([16]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return [16]byte{}, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return [16]byte{}, err
	}

	var result [16]byte
	copy(result[:], hash.Sum(nil))
	return result, nil
}

func computePHash(path string) (*goimagehash.ImageHash, error) {
	img, err := openImage(path)
	if err != nil {
		return nil, err
	}

	return goimagehash.PerceptionHash(img)
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
