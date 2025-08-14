# My Static Site

Built with [go-static](https://github.com/ahoglund/go-static) and deployed on GitHub Pages.

## Development

### Prerequisites

- Go 1.19 or later
- go-static installed: `go install github.com/ahoglund/go-static@latest`

### Local Development

```bash
# Build the site
go-static build

# Serve locally (http://localhost:8080)
go-static serve

# Build with verbose output
go-static build --verbose

# Serve on custom port
go-static serve --port 3000
```

### Project Structure

```
.
├── pages/          # Markdown and HTML source files
│   ├── index.md    # Homepage
│   └── about.md    # About page
├── templates/      # Go template files
│   ├── header.tmpl
│   ├── footer.tmpl
│   ├── nav.tmpl
│   ├── content.tmpl
│   └── index.tmpl
├── assets/         # Static assets
│   └── css/
│       └── main.css # Tailwind CSS
└── public/         # Generated output (ignored by git)
```

## Deployment

This site is automatically deployed to GitHub Pages using GitHub Actions.

### Setup GitHub Pages

1. Go to your repository settings
2. Navigate to "Pages" in the sidebar
3. Under "Source", select "GitHub Actions"
4. The workflow will automatically deploy on pushes to the `main` branch

### Custom Domain (Optional)

To use a custom domain:

1. Add a `CNAME` file to the `assets/` directory with your domain name
2. Configure your domain's DNS to point to GitHub Pages
3. Enable "Enforce HTTPS" in repository settings

## Adding Content

### New Page

Create a new `.md` file in the `pages/` directory:

```markdown
---
title: My New Page
---

# My New Page

Content goes here...
```

### Blog Post

Create a new file in `pages/blog/`:

```markdown
---
title: My First Post
template: post
date: 2024-01-15
---

# My First Post

Blog content here...
```

### Custom Template

Create a new template in `templates/`:

```html
{{define "post"}}
{{template "header" .}}
<article class="prose prose-lg mx-auto">
    <h1>{{.title}}</h1>
    <time class="text-gray-600">{{.date}}</time>
    {{.content}}
</article>
{{template "footer" .}}
{{end}}
```

## Styling

This site uses [Tailwind CSS](https://tailwindcss.com) for styling. Customize the design by:

1. Editing `assets/css/main.css`
2. Modifying the templates in `templates/`
3. Using Tailwind utility classes in your markdown

## License

[MIT License](LICENSE) - feel free to use this template for your own projects!