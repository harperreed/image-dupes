package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/vitali-fedulov/images4"
)

type ImageInfo struct {
	Path     string
	FileHash [16]byte
	Icon     images4.IconT
}

type ImageOpener interface {
	Open(path string) (image.Image, error)
}

type IconCreator interface {
	Icon(img image.Image) images4.IconT
}

type FileHasher interface {
	ComputeFileHash(path string) ([16]byte, error)
}

type DefaultImageOpener struct{}

func (d DefaultImageOpener) Open(path string) (image.Image, error) {
	return images4.Open(path)
}

type DefaultIconCreator struct{}

func (d DefaultIconCreator) Icon(img image.Image) images4.IconT {
	return images4.Icon(img)
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
	var wg sync.WaitGroup
	imageInfos := make([]ImageInfo, len(imagePaths))
	errors := make(chan error, len(imagePaths))

	for i, path := range imagePaths {
		wg.Add(1)
		go func(index int, imagePath string) {
			defer wg.Done()

			fileHash, err := hasher.ComputeFileHash(imagePath)
			if err != nil {
				errors <- fmt.Errorf("Error computing file hash for %s: %v", imagePath, err)
				progress <- fmt.Sprintf("Skipped %d/%d: %s (hash error)", index+1, len(imagePaths), filepath.Base(imagePath))
				return
			}

			img, err := opener.Open(imagePath)
			if err != nil {
				errors <- fmt.Errorf("Error opening image %s: %v", imagePath, err)
				progress <- fmt.Sprintf("Skipped %d/%d: %s (open error)", index+1, len(imagePaths), filepath.Base(imagePath))
				return
			}

			icon := iconCreator.Icon(img)

			imageInfos[index] = ImageInfo{Path: imagePath, FileHash: fileHash, Icon: icon}
			progress <- fmt.Sprintf("Processed %d/%d: %s", index+1, len(imagePaths), filepath.Base(imagePath))
		}(i, path)
	}

	wg.Wait()
	close(errors)

	var errSlice []error
	for err := range errors {
		errSlice = append(errSlice, err)
	}

	if len(errSlice) > 0 {
		return imageInfos, fmt.Errorf("Errors occurred during hash computation: %v", errSlice)
	}

	return imageInfos, nil
}
