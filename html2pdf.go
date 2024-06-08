package html2pdf

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type Generator[T any] struct {
	OutputPath     string             // directory path for generated data
	FinalPdf       string             // the merged pdf name (make sure to include .pdf in the name)
	Template       *template.Template // html template
	Data           []T                // valid data for feeding the template
	browser        *rod.Browser       // rod browser for auto generating pdf from html views
	HtmlFiles      []*os.File         // list of generated html files
	PdfFiles       []string           // list of generated pdf files
	SingleHtmlFile bool               // If you want the template to be single file only
	Linux          bool               // set this to true if running on linux to allow NoSandbox option on Rod luncher
}

// Generate pdf file from multible html templates
func (g *Generator[T]) CreatePdf() error {
	err := g.GenerateTemplates()
	if err != nil {
		return err
	}

	l := launcher.New().Headless(true).Leakless(true).NoSandbox(g.Linux)
	g.browser = rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer g.browser.MustClose()

	for i, file := range g.HtmlFiles {
		defer file.Close()
		pdfFilePath := fmt.Sprintf("./%s/output%d.pdf", g.OutputPath, i)

		cd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory path : %v", err.Error())
		}

		filePath := filepath.Join(cd, file.Name())
		filePath = strings.ReplaceAll(filePath, "\\", "/")

		// Construct the URL using the file path
		url := "file://" + filePath

		err = g.CapturePDF(g.browser, url, pdfFilePath)
		if err != nil {
			return fmt.Errorf("error capturing PDF files: %v", err.Error())
		}
		g.PdfFiles = append(g.PdfFiles, pdfFilePath)
	}

	if !g.SingleHtmlFile {
		// Merge the PDF files
		err = g.MergePDFs(g.PdfFiles, g.FinalPdf)
		if err != nil {
			return fmt.Errorf("error merging PDF files: %v", err.Error())
		}
	} else {
		output := fmt.Sprint("./" + g.FinalPdf)
		copyFile(g.PdfFiles[0], output)
	}

	return nil
}

// Generate html templates from the given data and save them into .g.OutputPath
func (g *Generator[T]) GenerateTemplates() error {
	if g.SingleHtmlFile {
		file, err := g.CreateHtmlFile(1)
		if err != nil {
			return err
		}
		err = g.Template.Execute(file, g.Data)
		if err != nil {
			return fmt.Errorf("error generating templates: %v", err.Error())
		}
		g.HtmlFiles = append(g.HtmlFiles, file)
		return nil
	}

	// for multible html files
	for i, v := range g.Data {
		file, err := g.CreateHtmlFile(i)
		if err != nil {
			return err
		}
		err = g.Template.Execute(file, v)
		if err != nil {
			return fmt.Errorf("error generating templates: %v", err.Error())
		}
		g.HtmlFiles = append(g.HtmlFiles, file)
	}

	return nil
}

func (g *Generator[T]) CreateHtmlFile(id int) (*os.File, error) {
	os.Mkdir(g.OutputPath, 0755)
	name := fmt.Sprintf("./%s/output%d.html", g.OutputPath, id)
	file, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("error creating html files: %v", err.Error())
	}
	return file, nil
}

// Delete html and pdf files except the merged pdf
func (g *Generator[T]) DeleteFiles() error {
	err := os.RemoveAll(g.OutputPath)
	if err != nil {
		return fmt.Errorf("error deleting files directory: %v", err)
	}
	return nil
}

// Automate opening a prowser then capture the html page as single pdf file
func (g *Generator[T]) CapturePDF(browser *rod.Browser, htmlUrl, outputPath string) error {
	page, err := browser.Page(proto.TargetCreateTarget{URL: htmlUrl})
	if err != nil {
		return fmt.Errorf("error creating browser page: %v", err)
	}
	page.MustWaitLoad().MustPDF(outputPath)
	fmt.Println(htmlUrl, ":::::", outputPath)
	return nil
}

// Merging all generated pdf together and create the output file
func (g *Generator[T]) MergePDFs(inputFiles []string, outputFile string) error {
	// Merge the PDF files
	err := api.MergeCreateFile(inputFiles, outputFile, false, api.LoadConfiguration())
	if err != nil {
		return fmt.Errorf("error merging PDF files: %v", err.Error())
	}
	return nil
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the content from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush file to disk
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
