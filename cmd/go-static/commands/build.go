package commands

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/ahoglund/go-static/pkg/config"
	"github.com/ahoglund/go-static/pkg/processor"
	"github.com/ahoglund/go-static/pkg/template"
	"github.com/spf13/cobra"
)

var (
	buildOutput string
	buildClean  bool
)

var buildCmd = &cobra.Command{
	Use:   "build [directory]",
	Short: "Build the static site",
	Long: `Build the static site from markdown and template files.

The directory should contain:
- pages/     - Markdown and HTML source files
- templates/ - Go template files
- assets/    - Static assets (optional)

Output will be generated in the public/ directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		if verbose {
			fmt.Printf("Building site from: %s\n", targetDir)
		}

		cfg := config.NewConfig(targetDir)
		if buildOutput != "" {
			cfg.PublicDir = buildOutput
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}

		if verbose {
			fmt.Printf("Templates: %s\n", cfg.TemplateDir)
			fmt.Printf("Pages: %s\n", cfg.PagesDir)
			fmt.Printf("Output: %s\n", cfg.PublicDir)
			fmt.Printf("Assets: %s\n", cfg.AssetsDir)
		}

		templateLoader := template.NewTemplateLoader(cfg)
		templates, err := templateLoader.LoadTemplates()
		if err != nil {
			return fmt.Errorf("template loading error: %w", err)
		}

		pageProcessor := processor.NewPageProcessor(cfg, templates)

		var processedFiles int
		err = filepath.WalkDir(cfg.PagesDir, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if verbose {
				fmt.Printf("Processing: %s\n", path)
			}
			
			processedFiles++
			return pageProcessor.ProcessPage(path)
		})

		if err != nil {
			return fmt.Errorf("page processing error: %w", err)
		}

		err = processor.ProcessAssets(cfg.AssetsDir, cfg.PublicDir, verbose)
		if err != nil {
			if verbose {
				log.Printf("Asset processing warning: %v", err)
			}
		}

		fmt.Printf("Site built successfully! Processed %d files.\n", processedFiles)
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "output directory (default: ./public)")
	buildCmd.Flags().BoolVar(&buildClean, "clean", false, "clean output directory before building")
}