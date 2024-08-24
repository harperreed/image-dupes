package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vitali-fedulov/images4"
)

type ImageInfo struct {
	Path     string
	FileHash [16]byte
	Icon     images4.IconT
}

func computeHashes(imagePaths []string, progress chan<- string) ([]ImageInfo, error) {
	var imageInfos []ImageInfo
	for i, path := range imagePaths {
		fileHash, err := computeFileHash(path)
		if err != nil {
			fmt.Printf("Error computing file hash for %s: %v\n", path, err)
			continue
		}

		img, err := images4.Open(path)
		if err != nil {
			fmt.Printf("Error opening image %s: %v\n", path, err)
			continue
		}

		icon := images4.Icon(img)

		imageInfos = append(imageInfos, ImageInfo{Path: path, FileHash: fileHash, Icon: icon})

		// Send progress update
		progress <- fmt.Sprintf("Processed %d/%d: %s", i+1, len(imagePaths), filepath.Base(path))
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
