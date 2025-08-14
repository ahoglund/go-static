# go-static

A fast, opinionated static site generator written in Go with built-in Tailwind CSS support.

## Features

- **Fast Builds**: Efficient static site generation
- **Tailwind CSS**: Professional styling out of the box
- **Go Templates**: Powerful templating system with partials
- **Asset Processing**: Automatic CSS bundling and minification  
- **Development Server**: Built-in server for local development
- **CLI Interface**: Modern command-line interface with Cobra
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Installation

### Method 1: Go Install (Recommended)

If you have Go 1.19+ installed:

```bash
go install github.com/ahoglund/go-static@latest
```

This will install `go-static` to your `$GOPATH/bin` directory.

### Method 2: Download Pre-built Binary

1. Go to the [Releases page](https://github.com/ahoglund/go-static/releases)
2. Download the appropriate binary for your platform
3. Extract and place in your PATH

### Method 3: Build from Source

```bash
git clone https://github.com/ahoglund/go-static.git
cd go-static
make build
# Binary will be in ./bin/go-static
```

### Method 4: Homebrew (Coming Soon)

```bash
# Coming soon
brew install go-static
```

## Quick Start

### Local Development

```bash
# Create a new site
go-static init my-site

# Build the site  
go-static build my-site

# Serve locally (http://localhost:8080)
go-static serve my-site
```

### GitHub Pages Deployment

```bash
# Create a site optimized for GitHub Pages
go-static init my-site --github-pages

# Follow the setup instructions to deploy
cd my-site
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/username/my-site.git
git push -u origin main

# Enable GitHub Pages in repository settings
```

## Commands

- `go-static init [directory]` - Initialize a new site with Tailwind CSS
- `go-static build [directory]` - Build the static site
- `go-static serve [directory]` - Serve the site locally
- `go-static version` - Show version information

### Flags

**Global:**
- `--verbose, -v` - Verbose output

**Init:**
- `--github-pages` - Setup GitHub Pages deployment (adds .gitignore, workflow, README)

**Build:**
- `--output, -o` - Custom output directory (default: ./public)

**Serve:**
- `--port, -p` - Custom port (default: 8080)
- `--host` - Custom host (default: localhost)

## Project Structure

go-static follows conventions for directory structure:

```
.
├── pages/          # Markdown and HTML source files
│   ├── index.md
│   └── about.md
├── templates/      # Go template files
│   ├── header.tmpl
│   ├── footer.tmpl
│   ├── nav.tmpl
│   ├── content.tmpl
│   └── index.tmpl
├── assets/         # Static assets (CSS, images, etc.)
│   └── css/
│       └── main.css
└── public/         # Generated output (created by build)
    ├── index.html
    ├── about.html
    └── css/
        └── main.css
```

## Content Files

### Markdown with Frontmatter

Create `.md` files in the `pages/` directory with YAML frontmatter:

```markdown
---
title: My Page Title
template: index
---

# My Content

This is **markdown** content that will be converted to HTML.

- List item 1
- List item 2

[Link to another page](/about.html)
```

### Supported Frontmatter Fields

- `title` (required): Page title
- `template` (optional): Template to use (defaults to "index")

## Templates

Templates use Go's `text/template` syntax with custom components:

```html
{{define "index"}}
{{template "header" .}}
{{template "nav" .}}
{{template "content" .}}
{{template "footer" .}}
{{end}}
```

### Available Variables

- `{{.title}}` - Page title from frontmatter
- `{{.content}}` - Processed markdown content
- Any custom frontmatter fields

## CSS and Styling

go-static includes **Tailwind CSS** by default:

- Sites are scaffolded with a complete Tailwind setup
- CSS files are automatically processed and minified
- Custom Tailwind components for typography
- CDN version included for immediate styling

## GitHub Pages Deployment

go-static includes first-class support for GitHub Pages deployment.

### Automatic Setup

```bash
go-static init my-site --github-pages
```

This creates:
- **`.gitignore`** - Ignores `public/` directory and common development files
- **`.github/workflows/deploy.yml`** - GitHub Actions workflow for automatic deployment
- **`README.md`** - Documentation for your repository

### Deployment Process

1. **Create Repository**: Push your site to a GitHub repository
2. **Enable Pages**: Go to repository Settings → Pages → Source: "GitHub Actions"
3. **Auto-Deploy**: Every push to `main` branch automatically rebuilds and deploys your site

### Workflow Features

- **Fast Builds**: Uses latest go-static version
- **Automatic Deployment**: No manual intervention needed
- **Branch Protection**: Only deploys from `main` branch
- **Artifact Caching**: Efficient build process

### Custom CSS

Add custom CSS to `assets/css/main.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Your custom styles */
.my-custom-class {
    @apply text-blue-600 hover:text-blue-800;
}
```

## Development

### Local Development

```bash
# Start development server
go-static serve my-site --port 3000

# Build with verbose output
go-static build my-site --verbose

# Build to custom directory
go-static build my-site --output dist
```

### Building

```bash
# Development build
make build

# Cross-platform builds
make build-all

# Install locally
make install
```

## Examples

### Basic Blog

```bash
go-static init my-blog
cd my-blog

# Add a blog post
cat > pages/first-post.md << EOF
---
title: My First Post
---

# Hello World

This is my first blog post!
EOF

# Build and serve
go-static build
go-static serve
```

### Custom Template

Create a custom template in `templates/post.tmpl`:

```html
{{define "post"}}
{{template "header" .}}
<article class="prose prose-lg mx-auto">
    <h1>{{.title}}</h1>
    <div class="text-gray-600">{{.date}}</div>
    {{.content}}
</article>
{{template "footer" .}}
{{end}}
```

Use in frontmatter:

```markdown
---
title: My Post
template: post
date: 2024-01-15
---

Content here...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Disclaimer

This project is in active development. Features and APIs may change. Please check the [releases page](https://github.com/ahoglund/go-static/releases) for stable versions.