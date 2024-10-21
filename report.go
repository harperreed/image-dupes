package main

import (
	"html/template"
	"os"
)

type ImageData struct {
	Path           string
	PerceptualHash uint64
}

type HTMLData struct {
	Groups [][]ImageData
}

func generateHTMLReport(similarGroups [][]ImageInfo, outputFile string) error {
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
        .hash { font-size: 0.8em; color: #666; margin-top: 5px; }
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
                <img src="file://{{.Path}}" alt="Similar Image">
                <div class="path">{{.Path}}</div>
                <div class="hash">Perceptual Hash: {{printf "%016x" .PerceptualHash}}</div>
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

	var groups [][]ImageData
	for _, group := range similarGroups {
		var imageGroup []ImageData
		for _, img := range group {
			imageGroup = append(imageGroup, ImageData{Path: img.Path, PerceptualHash: img.PerceptualHash})
		}
		groups = append(groups, imageGroup)
	}

	data := HTMLData{Groups: groups}
	return t.Execute(file, data)
}
