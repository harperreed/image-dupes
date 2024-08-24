package main

import (
	"path/filepath"
	"os"
)

func scanDirectoryRecursive(rootDir string) ([]string, error) {
	var images []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				images = append(images, path)
			}
		}
		return nil
	})
	return images, err
}
