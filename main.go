package main

import (
	"flag"
	"fmt"
	"os"
)

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

	fmt.Println("Scanning directory for images...")
	images, err := scanDirectoryRecursive(*rootDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}
	fmt.Printf("Found %d images\n", len(images))

	progress := &Progress{totalFiles: len(images)}

	fmt.Println("Computing image hashes...")
	go displayProgress(progress)

	imageHashes, err := computeHashes(images, progress)
	if err != nil {
		fmt.Printf("Error computing hashes: %v\n", err)
		return
	}

	fmt.Println("\nFinding similar images...")
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
	fmt.Println("Generating HTML report...")
	err = generateHTMLReport(similarImages, *outputHTML)
	if err != nil {
		fmt.Printf("Error generating HTML report: %v\n", err)
		return
	}
	fmt.Printf("HTML report generated: %s\n", *outputHTML)
}
