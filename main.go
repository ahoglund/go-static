package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

// iterate through every page
// and write out the html version
// to the public directory.
func main() {
	// Get all .md files in pagesDir
	dirs, err := ioutil.ReadDir(pagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for _, file := range dirs {
		content, err := ioutil.ReadFile(pagesDir + "/" + file.Name())
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
		// Need to ignore comments in front matter.

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

		templateContent, err := ioutil.ReadFile(templateDir + "/" + config["template"] + ".html.template")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		frontMatter := &FrontMatter{
			title:    config["title"],
			template: config["template"],
			content:  markdown.ToHTML([]byte(markdownContent), nil, nil),
		}

		parsedTemplateContent := parseTemplate(templateContent, frontMatter)
		fmt.Println(parsedTemplateContent)
	}

}

// parseTemplate takes in a template containing {{ }} variable strings and replaces them
// with either the templates refered to (by t:<string>) in the config map.
// content is a special variable.
// title is a special variable.
func parseTemplate(templateContent []byte, frontMatter *FrontMatter) []string {
	parsedTemplateContent := make([]string, 0)
	for _, line := range strings.Split(string(templateContent), "\n") {
		r := regexp.MustCompile(`.*{{\s*(.*)\s*}}.*`)
		if r.Match([]byte(line)) {
			found := r.FindAllStringSubmatch(line, -1)
			v := strings.Split(found[0][1], ":")

			// I need to not trim space so much!
			varName := strings.TrimSpace(v[0])
			switch varName {
			case "t":
				// load the template and then add it to the parsedTemplateContent
				templateName := strings.TrimSpace(v[1])
				templateContent, err := ioutil.ReadFile(templateDir + "/" + templateName + ".html.template")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					os.Exit(1)
				}
				subTemplateContent := parseTemplate(templateContent, frontMatter)
				parsedTemplateContent = append(parsedTemplateContent, subTemplateContent...)
			case "content":
				parsedTemplateContent = append(parsedTemplateContent, string(frontMatter.content))
			case "title":
				parsedTemplateContent = append(parsedTemplateContent, string(frontMatter.title))
			default:
				fmt.Fprintf(os.Stderr, "Unsupported variable: %s", v[0])
				os.Exit(1)
			}
		} else {
			parsedTemplateContent = append(parsedTemplateContent, line)
		}

		// if the line contains {{ something }}
		// then that means it has content to be replaced.
		// If not, the the line should just be added to the
		// parsed template output.
	}
	return parsedTemplateContent
}
