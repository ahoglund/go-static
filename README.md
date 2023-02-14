# go-static
Just another static site generator. Written in Go.

## Description

A simple command line application that will build out static html files from a combination of configuration files, html templates, and markdown content files.

## Usage

This program makes a few assumed conventions about the file structure of your website. Based on these conventions, go-static will parse and build static artifacts.

Give the following file structure:

```
.
├── pages
│   └── reviews
│       ├── reviews-01.md
│       └── reviews-02.md
└── templates
    ├── content.tmpl
    ├── footer.tmpl
    ├── header.tmpl
    ├── index.tmpl
    └── nav.tmpl
```

go-static will write out files to the `public` folder in the same directory as follows:

```
├── public
│   └── reviews
│       ├── reviews-01.html
│       └── reviews-02.html
```

If we inspect the contents of one of the review markdown files, we get an idea of how this happens:

```
---
title: Review 01
---

This is markdown and should be rendered into the `content` variable in the template.

## This is a header


**this is bold**


- this
- is
- a
- list
```

This is just a basic markdown with [frontmatter](https://markdoc.dev/docs/frontmatter) to provide configuration metadata. The `title` key is required. You can also pass it a `template` key to specify which template in the template directory you want this markdown file to use. If none is given, it defaults to `index`.


Lets take a look at a template file:

```
{{ define "index" }}
{{ template "header" . }}
{{ template "nav" . }}
{{ template "content" . }}
{{ template "footer" . }}
{{ end }}
```

The `template` denotes that this `{{ }}` section should replaced with a template with the name on the inside the quotes. In the first line we see it should render the `header` template.

Let's take a look at that one:

```
<!DOCTYPE html>
<html>
<head>
<title>{{ .title }}</title>
</head>
```

Pretty basic html, but with some [mustache style](https://github.com/janl/mustache.js/) variable. In this case, there is a special variable `title` that will be replaced with the `title:` that is provided in the frontmatter of the markdown file.

The other special variable in templates is `content`, which can be specified in the same way: `{{ .content }}`. This will render out the html version of the markdown.

You can compose templates as you like, and even nest them recursively in order to compose html pages from reusable fragments.

## Disclaimer

This project is in _very early_ development stage and has lots of features missing and problem some bugs as well. Feel free to open an issue/pull request if you'd like to contribute or make suggestions for improvement.

