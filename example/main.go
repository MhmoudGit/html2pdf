package main

import (
	"fmt"
	"html/template"
	"github.com/MhmoudGit/html2pdf"
	"time"
)

type Todo struct {
	ID   int
	Item string
}

func main() {
	// Create a temporary directory for testing
	tempDir := "files"
	mergedPdf := "example.pdf"

	todos := []Todo{
		{ID: 1, Item: "First item"},
		{ID: 2, Item: "Second item"},
		{ID: 3, Item: "Third item"},
		{ID: 4, Item: "Fourth item"},
		{ID: 5, Item: "Fifth item"},
	}

	templ, err := template.ParseFiles("example/example.html")
	if err != nil {
		fmt.Println(err)
	}

	start := time.Now()
	g := html2pdf.Generator[Todo]{
		OutputPath: tempDir,
		FinalPdf:   mergedPdf,
		Template:   templ,
		Data:       todos,
	}
	err = g.CreatePdf()
	if err != nil {
		fmt.Println(err)
	}

	err = g.DeleteFiles()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(time.Since(start))
}
