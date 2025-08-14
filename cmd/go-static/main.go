package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ahoglund/go-static/pkg/config"
	"github.com/ahoglund/go-static/pkg/processor"
	"github.com/ahoglund/go-static/pkg/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "No target directory provided.")
		os.Exit(1)
	}

	targetDir := os.Args[1]
	if targetDir == "" {
		fmt.Fprint(os.Stderr, "No target directory provided.")
		os.Exit(1)
	}

	cfg := config.NewConfig(targetDir)
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	templateLoader := template.NewTemplateLoader(cfg)
	templates, err := templateLoader.LoadTemplates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template loading error: %v\n", err)
		os.Exit(1)
	}

	pageProcessor := processor.NewPageProcessor(cfg, templates)

	err = filepath.WalkDir(cfg.PagesDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return pageProcessor.ProcessPage(path)
	})
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Page processing error: %v\n", err)
		os.Exit(1)
	}

	err = processor.ProcessAssets(cfg.AssetsDir, cfg.PublicDir)
	if err != nil {
		log.Printf("Asset processing error: %v", err)
	}

	fmt.Println("Site built successfully!")
}