package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGenerateHTMLReport(t *testing.T) {
	tests := []struct {
		name          string
		similarGroups [][]string
		outputFile    string
		wantErr       bool
		checkContent  func(t *testing.T, content string)
	}{
		{
			name: "Successful HTML report generation",
			similarGroups: [][]string{
				{"/path/to/image1.jpg", "/path/to/image2.jpg"},
				{"/path/to/image3.jpg", "/path/to/image4.jpg", "/path/to/image5.jpg"},
			},
			outputFile: "test_report.html",
			wantErr:    false,
			checkContent: func(t *testing.T, content string) {
				expectedStrings := []string{
					"<!DOCTYPE html>",
					"<title>Similar Images Report</title>",
					"<h1>Similar Images Report</h1>",
					"<h2>Group 1</h2>",
					"<img src=\"file:///path/to/image1.jpg\"",
					"<img src=\"file:///path/to/image2.jpg\"",
					"<h2>Group 2</h2>",
					"<img src=\"file:///path/to/image3.jpg\"",
					"<img src=\"file:///path/to/image4.jpg\"",
					"<img src=\"file:///path/to/image5.jpg\"",
				}
				for _, str := range expectedStrings {
					if !strings.Contains(content, str) {
						t.Errorf("Generated HTML does not contain expected string: %s", str)
					}
				}
			},
		},
		{
			name:          "Empty similar groups input",
			similarGroups: [][]string{},
			outputFile:    "empty_report.html",
			wantErr:       false,
			checkContent: func(t *testing.T, content string) {
				if !strings.Contains(content, "<h1>Similar Images Report</h1>") {
					t.Errorf("Generated HTML does not contain the report title")
				}
				if strings.Contains(content, "<h2>Group") {
					t.Errorf("Generated HTML contains group headers when it shouldn't")
				}
			},
		},
		{
			name: "Large number of image groups",
			similarGroups: func() [][]string {
				groups := make([][]string, 1000)
				for i := range groups {
					groups[i] = []string{"/path/to/image1.jpg", "/path/to/image2.jpg"}
				}
				return groups
			}(),
			outputFile: "large_report.html",
			wantErr:    false,
			checkContent: func(t *testing.T, content string) {
				if !strings.Contains(content, "<h2>Group 1000</h2>") {
					t.Errorf("Generated HTML does not contain the last group")
				}
			},
		},
		{
			name: "Image paths with special characters",
			similarGroups: [][]string{
				{"/path/to/image with spaces.jpg", "/path/to/image_with_underscore.jpg"},
				{"/path/to/image_with_パーセント%編码.jpg"},
			},
			outputFile: "special_chars_report.html",
			wantErr:    false,
			checkContent: func(t *testing.T, content string) {
				expectedPaths := []string{
					"/path/to/image with spaces.jpg",
					"/path/to/image_with_underscore.jpg",
					"/path/to/image_with_パーセント%編码.jpg",
				}
				for _, path := range expectedPaths {
					if !strings.Contains(content, path) {
						t.Errorf("Generated HTML does not contain path with special characters: %s", path)
					}
				}
			},
		},
		{
			name:          "Invalid output file path",
			similarGroups: [][]string{{"image1.jpg"}},
			outputFile:    "/invalid/path/report.html",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateHTMLReport(tt.similarGroups, tt.outputFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateHTMLReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				content, err := ioutil.ReadFile(tt.outputFile)
				if err != nil {
					t.Fatalf("Failed to read generated HTML file: %v", err)
				}

				tt.checkContent(t, string(content))

				// Clean up the generated file
				os.Remove(tt.outputFile)
			}
		})
	}
}
