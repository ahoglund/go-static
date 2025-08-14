package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/site/*
var siteTemplates embed.FS

type Scaffolder struct {
	verbose bool
}

func NewScaffolder(verbose bool) *Scaffolder {
	return &Scaffolder{verbose: verbose}
}

func (s *Scaffolder) CreateSite(targetDir string) error {
	if s.verbose {
		fmt.Printf("Creating site structure in: %s\n", targetDir)
	}

	return fs.WalkDir(siteTemplates, "templates/site", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, "templates/site/")
		if relativePath == "" || relativePath == path {
			return nil
		}

		targetPath := filepath.Join(targetDir, relativePath)

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			if s.verbose {
				fmt.Printf("Created directory: %s\n", targetPath)
			}
			return nil
		}

		content, err := siteTemplates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}

		if s.verbose {
			fmt.Printf("Created file: %s\n", targetPath)
		}

		return nil
	})
}

func (s *Scaffolder) CreateDirectories(targetDir string) error {
	directories := []string{
		filepath.Join(targetDir, "assets"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		if s.verbose {
			fmt.Printf("Created directory: %s\n", dir)
		}
	}

	return nil
}