package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new static site",
	Long: `Initialize a new static site with the basic directory structure 
and example files.

Creates:
- pages/     - Directory for markdown and HTML files
- templates/ - Directory for Go template files  
- assets/    - Directory for static assets
- Example templates and a sample page`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		if verbose {
			fmt.Printf("Initializing new site in: %s\n", targetDir)
		}

		directories := []string{
			filepath.Join(targetDir, "pages"),
			filepath.Join(targetDir, "templates"),
			filepath.Join(targetDir, "assets"),
		}

		for _, dir := range directories {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
			if verbose {
				fmt.Printf("Created: %s\n", dir)
			}
		}

		files := map[string]string{
			filepath.Join(targetDir, "templates", "header.tmpl"): `{{define "header"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
</head>
<body>
{{end}}`,
			filepath.Join(targetDir, "templates", "footer.tmpl"): `{{define "footer"}}
</body>
</html>
{{end}}`,
			filepath.Join(targetDir, "templates", "nav.tmpl"): `{{define "nav"}}
<nav>
    <ul>
        <li><a href="/">Home</a></li>
    </ul>
</nav>
{{end}}`,
			filepath.Join(targetDir, "templates", "content.tmpl"): `{{define "content"}}
<main>
    {{.content}}
</main>
{{end}}`,
			filepath.Join(targetDir, "templates", "index.tmpl"): `{{define "index"}}
{{template "header" .}}
{{template "nav" .}}
{{template "content" .}}
{{template "footer" .}}
{{end}}`,
			filepath.Join(targetDir, "pages", "index.md"): `---
title: Welcome to My Site
---

# Welcome

This is your new static site built with go-static!

## Getting Started

1. Edit this file: ` + "`pages/index.md`" + `
2. Add more pages to the ` + "`pages/`" + ` directory
3. Customize templates in ` + "`templates/`" + `
4. Build with: ` + "`go-static build`" + `
5. Serve with: ` + "`go-static serve`" + `

## Features

- Markdown support
- Go templating
- Fast builds
- Simple structure`,
		}

		for filePath, content := range files {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", filePath, err)
			}
			if verbose {
				fmt.Printf("Created: %s\n", filePath)
			}
		}

		fmt.Printf("Site initialized successfully in: %s\n", targetDir)
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", targetDir)
		fmt.Println("  go-static build")
		fmt.Println("  go-static serve")

		return nil
	},
}