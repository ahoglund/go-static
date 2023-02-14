package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	title    string
	template string
}

const frontMatterDelimiter = "---\n"
const defaultTemplate = "index"

type Config struct {
	templateDir string
	pagesDir    string
	publicDir   string
	assetsDir   string
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
		assetsDir:   targetDir + "/assets",
	}

	templateFiles := []string{}
	err := filepath.WalkDir(config.templateDir, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			// should recurse here
			return nil
		}

		templateFiles = append(templateFiles, config.templateDir+"/"+info.Name())

		return nil
	})

	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Printf("Error parsing template files %v", err)
		os.Exit(1)
	}

	// Get all .md files in pagesDir
	err = filepath.WalkDir(config.pagesDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		processPage(path, ts, config)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// move all the assets to the public directory
	// srcDir := config.assetsDir
	// dstDir := config.publicDir
	//
	// err = processAssets(srcDir, dstDir)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	os.Exit(1)
	// }
}

func processAssets(srcDir string, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// return nil
			// err := processAssets(info.Name(), dstDir)
			// if err != nil {
			// 	return err
			// }
		}

		relPath, _ := filepath.Rel(srcDir, path)
		dstPath := filepath.Join(dstDir, relPath)

		return copyFile(path, dstPath)
	})

	if err != nil {
		return err
	}

	return nil
}

func processPage(file string, ts *template.Template, config *Config) {
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

	var y map[interface{}]interface{}
	err = yaml.Unmarshal([]byte(rawFrontMatter), &y)
	if err != nil {
		fmt.Printf("Error parsing YAML in file %s: %v", file, err)
		return
	}

	// If template is not set, then default to index
	if _, ok := y["template"]; !ok {
		y["template"] = defaultTemplate
	}

	if _, ok := y["title"]; !ok {
		fmt.Printf("File %s doesn't contain a title: ", file)
		return
	}

	var parsedPageBuf bytes.Buffer
	switch filepath.Ext(file) {
	case ".html":
		fmt.Fprint(&parsedPageBuf, rawContent)
	case ".md":
		fmt.Fprint(&parsedPageBuf, string(markdown.ToHTML([]byte(rawContent), nil, nil)))
	case ".tmpl":
		parsedTemplate, err := template.New("page").Parse(string(rawContent[:]))
		if err != nil {
			fmt.Printf("Error parsing template %s: %s", file, err)
			return
		}
		err = parsedTemplate.Execute(&parsedPageBuf, y)
		if err != nil {
			fmt.Printf("Error executing template %s: %s", file, err)
			return
		}
	default:
		fmt.Printf("Unsupported file type")
		os.Exit(1)
	}

	y["content"] = parsedPageBuf.String()

	var parsedTemplateBuf bytes.Buffer
	err = ts.ExecuteTemplate(&parsedTemplateBuf, y["template"].(string), y)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// pages/publications/file.tmpl -> publications/file.tmpl
	fileName := strings.ReplaceAll(file, config.pagesDir, "")
	fmt.Println(parsedTemplateBuf.String())
	err = writeTemplate(fileName, parsedTemplateBuf.String(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func writeTemplate(name string, content string, config *Config) error {
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

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func readTemplate(name string) string {
	templateContent, err := os.ReadFile(name + ".tmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	return string(templateContent)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
