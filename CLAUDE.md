# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-static is a static site generator written in Go that converts markdown files with frontmatter into HTML pages using Go templates. The architecture follows a simple pipeline: pages (markdown) + templates → static HTML files.

## Build and Run Commands

- **Build**: `go build -o go-static ./cmd/go-static`
- **Commands Available**:
  - `./go-static init [directory]` - Initialize new site with templates
  - `./go-static build [directory]` - Build static site (default: current directory)
  - `./go-static serve [directory]` - Serve built site locally
  - `./go-static version` - Show version information
- **Flags**: 
  - `--verbose, -v` - Verbose output (global)
  - `--output, -o` - Custom output directory (build)
  - `--port, -p` - Custom port (serve, default: 8080)
- **Example**: `./go-static build example/` (builds the example site)

## Architecture

### Core Components

- **cmd/go-static/main.go**: Application entry point with orchestration logic
- **pkg/config/**: Configuration management with validation
- **pkg/processor/**: Page processing logic and asset handling
- **pkg/template/**: Template loading and management
- **FrontMatter struct**: Handles YAML frontmatter parsing (title, template fields)
- **PageProcessor**: Processes individual markdown/HTML/template files
- **ProcessAssets()**: Copies static assets from assets/ to public/

### File Processing Pipeline

1. Parse all template files from `templates/` directory using `text/template`
2. Walk through `pages/` directory processing each file:
   - Extract YAML frontmatter (title, template)
   - Convert markdown to HTML using `github.com/gomarkdown/markdown`
   - Execute template with frontmatter data and content
   - Write HTML output to `public/` directory
3. Process assets from `assets/` to `public/`:
   - CSS files are processed with esbuild (bundling, minification)
   - Tailwind CSS directives are processed
   - Other assets are copied preserving directory structure

### File Type Support

- `.md`: Markdown files with frontmatter
- `.html`: Raw HTML files with frontmatter  
- `.tmpl`: Template files that get processed with frontmatter variables

### Template System

Uses Go's `text/template` package with template composition:
- Templates are defined with `{{define "template-name"}}`
- Variable substitution: `{{.title}}`, `{{.content}}`
- Template inclusion: `{{template "header" .}}`
- Default template is "index" if not specified in frontmatter

## Dependencies

- `github.com/gomarkdown/markdown`: Markdown to HTML conversion
- `gopkg.in/yaml.v3`: YAML frontmatter parsing
- `github.com/spf13/cobra`: CLI framework
- `github.com/evanw/esbuild`: CSS processing and minification
- Go 1.19+ required

## CSS Processing

- **Tailwind CSS Support**: Sites are scaffolded with Tailwind CSS by default
- **Automatic Processing**: CSS files are automatically bundled and minified
- **Modern Approach**: Uses esbuild for fast CSS processing
- **Custom Components**: Includes custom Tailwind components for typography

## Directory Structure Conventions

The application expects this structure in the target directory:
```
target-directory/
├── pages/          # Markdown/HTML source files
├── templates/      # Go template files
├── assets/         # Static assets (optional)
└── public/         # Generated output (created automatically)
```