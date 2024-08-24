package main

import (
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math/bits"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

const (
	hashSize = 8
)

type ImageHash struct {
	Path string
	Hash uint64
}

type HTMLData struct {
	Groups [][]string
}

func main() {
	rootDir := flag.String("dir", "", "Root directory to scan for images")
	threshold := flag.Int("threshold", 10, "Maximum hamming distance to consider images similar")
	outputHTML := flag.String("output", "report.html", "Output HTML file name")
	flag.Parse()

	if *rootDir == "" {
		fmt.Println("Please specify a root directory using -dir flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	images, err := scanDirectoryRecursive(*rootDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	imageHashes, err := computeHashes(images)
	if err != nil {
		fmt.Printf("Error computing hashes: %v\n", err)
		return
	}

	similarImages := findSimilarImages(imageHashes, *threshold)

	// Print to console
	fmt.Println("Similar images:")
	for i, group := range similarImages {
		fmt.Printf("Group %d:\n", i+1)
		for _, img := range group {
			fmt.Printf("  %s\n", img)
		}
		fmt.Println()
	}

	// Generate HTML report
	err = generateHTMLReport(similarImages, *outputHTML)
	if err != nil {
		fmt.Printf("Error generating HTML report: %v\n", err)
		return
	}
	fmt.Printf("HTML report generated: %s\n", *outputHTML)
}

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

func computeHashes(images []string) ([]ImageHash, error) {
	var imageHashes []ImageHash
	for _, path := range images {
		hash, err := computeDHash(path)
		if err != nil {
			fmt.Printf("Error computing hash for %s: %v\n", path, err)
			continue
		}
		imageHashes = append(imageHashes, ImageHash{Path: path, Hash: hash})
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

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func generateHTMLReport(similarImages [][]string, outputFile string) error {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Similar Images Report</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; }
        h1 { color: #333; }
        .group { margin-bottom: 40px; border: 1px solid #ccc; padding: 20px; }
        .group h2 { color: #666; }
        .images { display: flex; flex-wrap: wrap; gap: 10px; }
        .image-container { max-width: 200px; }
        img { max-width: 100%; height: auto; border: 1px solid #ddd; }
        .path { font-size: 0.8em; word-break: break-all; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>Similar Images Report</h1>
    {{range $index, $group := .Groups}}
    <div class="group">
        <h2>Group {{add $index 1}}</h2>
        <div class="images">
            {{range $group}}
            <div class="image-container">
                <img src="file://{{.}}" alt="Similar Image">
                <div class="path">{{.}}</div>
            </div>
            {{end}}
        </div>
    </div>
    {{end}}
</body>
</html>
`

	t, err := template.New("report").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := HTMLData{Groups: similarImages}
	return t.Execute(file, data)
}
