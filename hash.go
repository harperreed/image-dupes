package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
	"github.com/vitali-fedulov/images4"
)

type ImageInfo struct {
	Path           string
	FileHash       [16]byte
	PerceptualHash uint64
	Icon           images4.IconT
}

type ImageOpener interface {
	Open(path string) (image.Image, error)
}

type IconCreator interface {
	Icon(img image.Image) (images4.IconT, uint64, error)
}

type FileHasher interface {
	ComputeFileHash(path string) ([16]byte, error)
}

type DefaultImageOpener struct{}

func (d DefaultImageOpener) Open(path string) (image.Image, error) {
	return images4.Open(path)
}

type DefaultIconCreator struct{}

func (d DefaultIconCreator) Icon(img image.Image) (images4.IconT, uint64, error) {
	icon := images4.Icon(img)
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return images4.IconT{}, 0, err
	}
	return icon, hash.GetHash(), nil
}

type DefaultFileHasher struct{}

func (d DefaultFileHasher) ComputeFileHash(path string) ([16]byte, error) {
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

func computeHashes(imagePaths []string, progress chan<- string, opener ImageOpener, iconCreator IconCreator, hasher FileHasher) ([]ImageInfo, error) {
	var imageInfos []ImageInfo
	for i, path := range imagePaths {
		fileHash, err := hasher.ComputeFileHash(path)
		if err != nil {
			fmt.Printf("Error computing file hash for %s: %v\n", path, err)
			progress <- fmt.Sprintf("Skipped %d/%d: %s (hash error)", i+1, len(imagePaths), filepath.Base(path))
			continue
		}

		img, err := opener.Open(path)
		if err != nil {
			fmt.Printf("Error opening image %s: %v\n", path, err)
			progress <- fmt.Sprintf("Skipped %d/%d: %s (open error)", i+1, len(imagePaths), filepath.Base(path))
			continue
		}

		icon, perceptualHash, err := iconCreator.Icon(img)
		if err != nil {
			fmt.Printf("Error creating icon and perceptual hash for %s: %v\n", path, err)
			progress <- fmt.Sprintf("Skipped %d/%d: %s (icon/hash error)", i+1, len(imagePaths), filepath.Base(path))
			continue
		}

		imageInfos = append(imageInfos, ImageInfo{Path: path, FileHash: fileHash, PerceptualHash: perceptualHash, Icon: icon})

		progress <- fmt.Sprintf("Processed %d/%d: %s", i+1, len(imagePaths), filepath.Base(path))
	}
	return imageInfos, nil
}
