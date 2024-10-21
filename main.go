package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	rootDir := flag.String("dir", "", "Root directory to scan for images")
	outputHTML := flag.String("output", "report.html", "Output HTML file name")
	flag.Parse()

	if *rootDir == "" {
		fmt.Println("Please specify a root directory using -dir flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Scanning directory
	fmt.Println("Scanning directory for images...")
	images, err := scanDirectoryRecursive(*rootDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}
	fmt.Printf("Found %d images\n", len(images))

	// Computing hashes
	fmt.Println("Computing image hashes...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	progressChan := make(chan string)
	go func() {
		for msg := range progressChan {
			s.Suffix = " " + msg
		}
	}()

	imageInfos, err := computeHashes(images, progressChan, DefaultImageOpener{}, DefaultIconCreator{}, DefaultFileHasher{})
	s.Stop()
	close(progressChan)

	if err != nil {
		fmt.Printf("Error computing hashes: %v\n", err)
		return
	}
	fmt.Printf("Computed hashes for %d images\n", len(imageInfos))

	// Finding similar images
	fmt.Println("Finding similar images...")
	similarGroups := findSimilarImages(imageInfos)
	fmt.Printf("Found %d groups of similar images\n", len(similarGroups))

	// Generating HTML report
	fmt.Println("Generating HTML report...")
	s.Suffix = ""
	s.Start()
	err = generateHTMLReport(similarGroups, *outputHTML)
	s.Stop()
	if err != nil {
		fmt.Printf("Error generating HTML report: %v\n", err)
		return
	}
	fmt.Printf("HTML report generated: %s\n", *outputHTML)
}
