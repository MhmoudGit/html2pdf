package html2pdf

import (
	"html/template"
	"os"
	"testing"
	"time"
)

type Todo struct {
	ID   int
	Item string
}

var (
	tmpl = `
<!DOCTYPE html>
<html>
<head>
<title>Todo Information</title>
</head>
<body>
<h1>{{.ID}}</h1>
<p>{{.Item}}.</p>
</body>
</html>
`

	// for single file test
	singleTmpl = `
<!DOCTYPE html>
<html>
<head>
<title>Todo Information</title>
</head>
<body>
{{range .}}
<h1>{{.ID}}</h1>
<p>{{.Item}}.</p>
{{end}}
</body>
</html>
`
	// Create a temporary directory for testing
	tempDir   = "files"
	mergedPdf = "final.pdf"
	g         = Generator[Todo]{
		OutputPath:     tempDir,
		FinalPdf:       mergedPdf,
		SingleHtmlFile: true, // test for single html file only
	}
)

func TestGenerator(t *testing.T) {
	expectedTodos := []Todo{
		{ID: 1, Item: "First item"},
		{ID: 2, Item: "Second item"},
		{ID: 3, Item: "Third item"},
		{ID: 4, Item: "Fourth item"},
	}

	if g.SingleHtmlFile {
		template, err := template.New("todo").Parse(singleTmpl)
		if err != nil {
			panic(err)
		}
		g.Template = template
	} else {
		template, err := template.New("todo").Parse(tmpl)
		if err != nil {
			panic(err)
		}
		g.Template = template
	}

	g.Data = expectedTodos

	result := g

	// Test if the Data field matches the expectedTodos
	if len(result.Data) != len(expectedTodos) {
		t.Errorf("Expected Data length %d, got %d", len(expectedTodos), len(result.Data))
		return
	}

	for i, expected := range expectedTodos {
		if result.Data[i] != expected {
			t.Errorf("Expected todo at index %d to be %v, got %v", i, expected, result.Data[i])
		}
	}

}

func TestCreatePdf(t *testing.T) {
	todos := []Todo{
		{ID: 1, Item: "First item"},
		{ID: 2, Item: "Second item"},
		{ID: 3, Item: "Third item"},
		{ID: 4, Item: "Fourth item"},
	}

	if g.SingleHtmlFile {
		template, err := template.New("todo").Parse(singleTmpl)
		if err != nil {
			panic(err)
		}
		g.Template = template
	} else {
		template, err := template.New("todo").Parse(tmpl)
		if err != nil {
			panic(err)
		}
		g.Template = template
	}
	g.Data = todos

	start := time.Now()

	// test 1
	err := g.CreatePdf()
	if err != nil {
		t.Errorf("Error creating PDF: %v", err)
	}

	t.Log(time.Since(start))
}

func TestDeletingDataFolder(t *testing.T) {
	err := g.DeleteFiles()
	if err != nil {
		t.Errorf("Error deleting data folder: %v", err)
	}

	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("Error folder is not deleted")
	}
}
