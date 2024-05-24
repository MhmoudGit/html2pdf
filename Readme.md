# html2pdf

html2pdf is a Go package designed specifically for converting HTML content to PDF documents. It wraps the functionality of [Rod]("https://github.com/go-rod/rod") for browser automation and [pdfcpu]("https://github.com/pdfcpu/pdfcpu") for PDF generation, providing a simple and efficient way to create PDFs from HTML content.

- this is very simple code, basically a wrapper, for my personal usage as i struggled to find easy go library to suit my needs for converting html to pdf

- it is not a perfect solution and not yet completed and probably have bugs :)

## Installation

`go get github.com/MhmoudGit/html2pdf`

## Usage

```go
package main

import (
	"fmt"
	"html/template"
	"github.com/MhmoudGit/html2pdf"
)

type Todo struct {
	ID   int
	Item string
}

var tmpl = `
<!DOCTYPE html>
<html>
	<head>
		<title>Todo Information</title>
	</head>
	<body>
			<h1 style="color: red">{{.ID}}</h1>
			<p style="color: green; font-weight: bold">{{.Item}}.</p>
	</body>
</html>
`

func main() {
    // directory for saving generated data
	tempDir := "files"
    // the final output
	mergedPdf := "example.pdf"

	todos := []Todo{
		{ID: 1, Item: "First item"},
	}

	template, err := template.New("todo").Parse(tmpl)
	if err != nil {
		fmt.Println(err)
	}

    // create a pdf generator
	g := html2pdf.Generator[Todo]{
		OutputPath: tempDir,
		FinalPdf:   mergedPdf,
		Template:   template,
		Data:       todos,
	}

    // generate pdf
	err = g.CreatePdf()
	if err != nil {
		fmt.Println(err)
	}

    // (optional)
    // delete the generated templates and pdf
	err = g.DeleteFiles()
	if err != nil {
		fmt.Println(err)
	}
}

```

- `Find example using html file and css styling in the example folder`

## Acknowledgments

- Rod - A powerful Go package for browser automation.
- PDFCPU - A PDF processing library written in Go.
