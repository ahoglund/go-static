package assets

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type CSSProcessor struct {
	inputDir  string
	outputDir string
	verbose   bool
}

func NewCSSProcessor(inputDir, outputDir string, verbose bool) *CSSProcessor {
	return &CSSProcessor{
		inputDir:  inputDir,
		outputDir: outputDir,
		verbose:   verbose,
	}
}

func (p *CSSProcessor) ProcessCSS() error {
	cssFiles := []string{}
	
	err := filepath.WalkDir(p.inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		if strings.HasSuffix(path, ".css") {
			cssFiles = append(cssFiles, path)
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to walk CSS directory: %w", err)
	}
	
	if len(cssFiles) == 0 {
		if p.verbose {
			fmt.Println("No CSS files found to process")
		}
		return nil
	}
	
	for _, cssFile := range cssFiles {
		if err := p.processFile(cssFile); err != nil {
			return fmt.Errorf("failed to process CSS file %s: %w", cssFile, err)
		}
	}
	
	return nil
}

func (p *CSSProcessor) processFile(inputPath string) error {
	relPath, err := filepath.Rel(p.inputDir, inputPath)
	if err != nil {
		return err
	}
	
	outputPath := filepath.Join(p.outputDir, relPath)
	
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	
	if p.verbose {
		fmt.Printf("Processing CSS: %s -> %s\n", inputPath, outputPath)
	}
	
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{inputPath},
		Outfile:     outputPath,
		Bundle:      true,
		MinifyWhitespace: true,
		MinifyIdentifiers: true,
		MinifySyntax: true,
		Write:       true,
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
		},
		Target: api.ES2020,
	})
	
	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "CSS Error: %s\n", err.Text)
		}
		return fmt.Errorf("CSS processing failed")
	}
	
	if len(result.Warnings) > 0 && p.verbose {
		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "CSS Warning: %s\n", warning.Text)
		}
	}
	
	return nil
}

func (p *CSSProcessor) ProcessTailwind(inputPath, outputPath string) error {
	if p.verbose {
		fmt.Printf("Processing Tailwind CSS: %s -> %s\n", inputPath, outputPath)
	}
	
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{inputPath},
		Outfile:     outputPath,
		Bundle:      true,
		MinifyWhitespace: true,
		MinifyIdentifiers: true,
		MinifySyntax: true,
		Write:       true,
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
		},
		Target: api.ES2020,
	})
	
	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "Tailwind CSS Error: %s\n", err.Text)
		}
		return fmt.Errorf("Tailwind CSS processing failed")
	}
	
	if len(result.Warnings) > 0 && p.verbose {
		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "Tailwind CSS Warning: %s\n", warning.Text)
		}
	}
	
	return nil
}