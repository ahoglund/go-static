package processor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ahoglund/go-static/pkg/assets"
	"github.com/ahoglund/go-static/pkg/config"
	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

type PageProcessor struct {
	config    *config.Config
	templates *template.Template
}

func NewPageProcessor(cfg *config.Config, templates *template.Template) *PageProcessor {
	return &PageProcessor{
		config:    cfg,
		templates: templates,
	}
}

func (p *PageProcessor) ProcessPage(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file, err)
	}

	data := strings.Split(string(content), FrontMatterDelimiter)
	if len(data) < 3 {
		return fmt.Errorf("invalid frontmatter format in file %s", file)
	}

	rawFrontMatter := data[1]
	rawContent := data[2]

	var y map[interface{}]interface{}
	err = yaml.Unmarshal([]byte(rawFrontMatter), &y)
	if err != nil {
		return fmt.Errorf("error parsing YAML in file %s: %w", file, err)
	}

	if _, ok := y["template"]; !ok {
		y["template"] = DefaultTemplate
	}

	if _, ok := y["title"]; !ok {
		return fmt.Errorf("file %s doesn't contain a title", file)
	}

	var parsedPageBuf bytes.Buffer
	switch filepath.Ext(file) {
	case ".html":
		fmt.Fprint(&parsedPageBuf, rawContent)
	case ".md":
		fmt.Fprint(&parsedPageBuf, string(markdown.ToHTML([]byte(rawContent), nil, nil)))
	case ".tmpl":
		parsedTemplate, err := template.New(file).Parse(rawContent)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %w", file, err)
		}
		err = parsedTemplate.Execute(&parsedPageBuf, y)
		if err != nil {
			return fmt.Errorf("error executing template %s: %w", file, err)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(file))
	}

	y["content"] = parsedPageBuf.String()

	var parsedTemplateBuf bytes.Buffer
	err = p.templates.ExecuteTemplate(&parsedTemplateBuf, y["template"].(string), y)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	relativePath, err := filepath.Rel(p.config.PagesDir, file)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}
	err = p.writeTemplate(relativePath, parsedTemplateBuf.String())
	if err != nil {
		return fmt.Errorf("error writing template: %w", err)
	}

	return nil
}

func (p *PageProcessor) writeTemplate(name string, content string) error {
	err := os.MkdirAll(p.config.PublicDir+"/"+filepath.Dir(name), os.ModePerm)
	if err != nil {
		return err
	}

	newFileName := strings.Replace(name, filepath.Ext(name), ".html", 1)
	file, err := os.Create(p.config.PublicDir + "/" + newFileName)
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

func ProcessAssets(srcDir string, dstDir string, verbose bool) error {
	cssProcessor := assets.NewCSSProcessor(srcDir, dstDir, verbose)
	
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		if strings.HasSuffix(path, ".css") {
			return cssProcessor.ProcessTailwind(path, dstPath)
		}

		return copyFile(path, dstPath)
	})

	return err
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
