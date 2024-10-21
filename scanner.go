package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func scanDirectoryRecursive(rootDir string) ([]string, error) {
	// First, check if the rootDir is actually a directory
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("the provided path is not a directory: %s", rootDir)
	}

	var images []string
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			lowerExt := strings.ToLower(ext)
			if lowerExt == ".jpg" || lowerExt == ".jpeg" || lowerExt == ".png" {
				images = append(images, path)
			}
		}
		return nil
	})
	return images, err
}
