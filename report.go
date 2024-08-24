package main

import (
	"html/template"
	"os"
)

type HTMLData struct {
	Groups [][]string
}

func generateHTMLReport(similarGroups [][]string, outputFile string) error {
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

	data := HTMLData{Groups: similarGroups}
	return t.Execute(file, data)
}
