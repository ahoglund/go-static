package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	title    string
	template string
	content  *bytes.Buffer
}

const frontMatterDelimiter = "---\n"
const defaultTemplate = "index"

type Config struct {
	templateDir string
	pagesDir    string
	publicDir   string
}

// iterate through every page
// and write out the html version
// to the public directory.
func main() {
	targetDir := os.Args[1]

	if targetDir == "" {
		fmt.Fprint(os.Stderr, "No target directory provided.")
		os.Exit(1)
	}

	// strip trailing slash from target Dir
	targetDir = strings.TrimSuffix(targetDir, "/")
	config := &Config{
		templateDir: targetDir + "/templates",
		pagesDir:    targetDir + "/pages",
		publicDir:   targetDir + "/public",
	}

	// Get all .md files in pagesDir
	err := filepath.WalkDir(config.pagesDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		processTemplate(path, config)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func processTemplate(file string, config *Config) {
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
	rawContent := data[2]

	// TODO: Need to ignore comments in front matter.

	frontmatter := make(map[string]string)

	for _, line := range strings.Split((rawFrontMatter), "\n") {
		if line == "" {
			continue
		}
		d := strings.Split(line, ":")
		frontmatter[strings.TrimSpace(d[0])] = strings.TrimSpace(d[1])
	}

	// If template is not set, then default to index
	if frontmatter["template"] == "" {
		frontmatter["template"] = defaultTemplate
	}

	// detect file type. process accordingly.
	var parsedContent bytes.Buffer
	switch filepath.Ext(file) {
	case ".html":
		fmt.Fprint(&parsedContent, rawContent)
	case ".md":
		fmt.Fprint(&parsedContent, string(markdown.ToHTML([]byte(rawContent), nil, nil)))
	case ".tmpl":
		parsedTemplate, _ := template.New("foo").Parse(string(rawContent[:]))
		var m map[interface{}]interface{}
		err := yaml.Unmarshal([]byte(rawFrontMatter), &m)
		if err != nil {
			fmt.Printf("Error parsing YAML: %v", err)
			return
		}

		parsedTemplate.Execute(&parsedContent, m)
	default:

	}

	frontMatter := &FrontMatter{
		title:    frontmatter["title"],
		template: config.templateDir + "/" + frontmatter["template"],
		content:  &parsedContent,
	}

	fileName := strings.ReplaceAll(file, config.pagesDir, "")
	parsedTemplateContent := parseTemplate(readTemplate(frontMatter.template), frontMatter, config)
	err = writeTemplate(fileName, renderTemplate(parsedTemplateContent), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func writeTemplate(name string, content []string, config *Config) error {
	err := os.MkdirAll(config.publicDir+"/"+filepath.Dir(name), os.ModePerm)
	if err != nil {
		return err
	}

	newFileName := strings.Replace(name, filepath.Ext(name), ".html", 1)
	file, err := os.Create(config.publicDir + "/" + newFileName)
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
	templateContent, err := os.ReadFile(name + ".html.template")
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
func parseTemplate(templateContent []byte, frontMatter *FrontMatter, config *Config) []string {
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
				templateName := config.templateDir + "/" + varContent[1]
				subTemplateContent := make([]string, 0)
				subTemplateContent = append(subTemplateContent, beforeContent)
				subTemplateContent = append(subTemplateContent, parseTemplate(readTemplate(templateName), frontMatter, config)...)
				subTemplateContent = append(subTemplateContent, afterContent)

				parsedTemplateContent = append(parsedTemplateContent, subTemplateContent...)
			case "content":
				parsedTemplateContent = append(parsedTemplateContent, beforeContent+string(frontMatter.content.Bytes())+afterContent)
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
