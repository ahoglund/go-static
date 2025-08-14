package config

import (
	"fmt"
	"strings"
)

type Config struct {
	TemplateDir string
	PagesDir    string
	PublicDir   string
	AssetsDir   string
}

func NewConfig(targetDir string) *Config {
	targetDir = strings.TrimSuffix(targetDir, "/")
	
	return &Config{
		TemplateDir: targetDir + "/templates",
		PagesDir:    targetDir + "/pages",
		PublicDir:   targetDir + "/public",
		AssetsDir:   targetDir + "/assets",
	}
}

func (c *Config) Validate() error {
	if c.TemplateDir == "" {
		return fmt.Errorf("template directory cannot be empty")
	}
	if c.PagesDir == "" {
		return fmt.Errorf("pages directory cannot be empty")
	}
	if c.PublicDir == "" {
		return fmt.Errorf("public directory cannot be empty")
	}
	return nil
}