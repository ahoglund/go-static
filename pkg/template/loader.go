package template

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/ahoglund/go-static/pkg/config"
)

type TemplateLoader struct {
	config *config.Config
}

func NewTemplateLoader(cfg *config.Config) *TemplateLoader {
	return &TemplateLoader{
		config: cfg,
	}
}

func (t *TemplateLoader) LoadTemplates() (*template.Template, error) {
	templateFiles := []string{}
	
	err := filepath.WalkDir(t.config.TemplateDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}

		templateFiles = append(templateFiles, t.config.TemplateDir+"/"+info.Name())
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk template directory: %w", err)
	}

	if len(templateFiles) == 0 {
		return nil, fmt.Errorf("no template files found in %s", t.config.TemplateDir)
	}

	templates, err := template.ParseFiles(templateFiles...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template files: %w", err)
	}

	return templates, nil
}
