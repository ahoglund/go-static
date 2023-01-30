package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
)

type FrontMatter struct {
	title    string
	template string
	content  []byte
}

const templateDir = "example/templates"
const pagesDir = "example/pages"
const frontMatterDelimiter = "---\n"
const defaultTemplate = "index"
const publicDir = "example/public"

// iterate through every page
// and write out the html version
// to the public directory.
func main() {
	// Get all .md files in pagesDir
	err := filepath.WalkDir(pagesDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dealWithTemplate(path)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func dealWithTemplate(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	data := make([]string, 0)
	for _, line := range strings.Split(string(content), frontMatterDelimiter) {
		data = append(data, line)
	}

	rawFrontMatter := data[1]
	markdownContent := data[2]

	// TODO: Need to ignore comments in front matter.

	config := make(map[string]string)

	for _, line := range strings.Split((rawFrontMatter), "\n") {
		if line == "" {
			continue
		}
		d := strings.Split(line, ":")
		config[strings.TrimSpace(d[0])] = strings.TrimSpace(d[1])
	}

	// If template is not set, then default to index
	if config["template"] == "" {
		config["template"] = defaultTemplate
	}

	frontMatter := &FrontMatter{
		title:    config["title"],
		template: config["template"],
		content:  markdown.ToHTML([]byte(markdownContent), nil, nil),
	}

	fileName := strings.ReplaceAll(file, pagesDir, "")
	fileName = strings.ReplaceAll(fileName, ".md", "")
	parsedTemplateContent := parseTemplate(readTemplate(frontMatter.template), frontMatter)
	err = writeTemplate(fileName, renderTemplate(parsedTemplateContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func writeTemplate(name string, content []string) error {
	err := os.MkdirAll(publicDir+"/"+filepath.Dir(name), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(publicDir + "/" + name + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each string to the file
	for _, line := range content {
		_, err = file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderTemplate(content []string) []string {
	finalContent := make([]string, 0)

	for _, line := range content {
		if line == "" {
			continue
		} else {
			finalContent = append(finalContent, line)
		}
		finalContent = append(finalContent, "\n")
	}
	return finalContent
}

func readTemplate(name string) []byte {
	templateContent, err := os.ReadFile(templateDir + "/" + name + ".html.template")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	return templateContent
}

// parseTemplate takes in a template containing {{ }} variable strings and replaces them
// with either the templates refered to (by t:<string>) in the config map.
// content is a special variable.
// title is a special variable.
func parseTemplate(templateContent []byte, frontMatter *FrontMatter) []string {
	parsedTemplateContent := make([]string, 0)
	for _, line := range strings.Split(string(templateContent), "\n") {
		r := regexp.MustCompile(`(.*){{\s*([a-zA-Z:]+)\s*}}(.*)`)
		if r.Match([]byte(line)) {
			found := r.FindAllStringSubmatch(line, -1)
			// I need to not trim space so much!
			beforeContent := found[0][1]
			varContent := strings.Split(found[0][2], ":")
			afterContent := found[0][3]
			varName := varContent[0]
			switch varName {
			case "t":
				templateName := varContent[1]
				subTemplateContent := make([]string, 0)
				subTemplateContent = append(subTemplateContent, beforeContent)
				subTemplateContent = append(subTemplateContent, parseTemplate(readTemplate(templateName), frontMatter)...)
				subTemplateContent = append(subTemplateContent, afterContent)

				parsedTemplateContent = append(parsedTemplateContent, subTemplateContent...)
			case "content":
				parsedTemplateContent = append(parsedTemplateContent, beforeContent+string(frontMatter.content)+afterContent)
			case "title":
				parsedTemplateContent = append(parsedTemplateContent, beforeContent+string(frontMatter.title)+afterContent)
			default:
				fmt.Fprintf(os.Stderr, "Unsupported variable: %s", varName)
				os.Exit(1)
			}
		} else {
			parsedTemplateContent = append(parsedTemplateContent, line)
		}
	}
	return parsedTemplateContent
}
