package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahoglund/go-static/pkg/config"
	"github.com/ahoglund/go-static/pkg/processor"
	"github.com/ahoglund/go-static/pkg/template"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	servePort string
	serveHost string
)

func buildSite(targetDir string) error {
	cfg := config.NewConfig(targetDir)
	
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	templateLoader := template.NewTemplateLoader(cfg)
	templates, err := templateLoader.LoadTemplates()
	if err != nil {
		return fmt.Errorf("template loading error: %w", err)
	}

	pageProcessor := processor.NewPageProcessor(cfg, templates)

	err = filepath.WalkDir(cfg.PagesDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if strings.Contains(fileName, ".tmp") || strings.HasSuffix(fileName, ".swp") || 
		   strings.HasPrefix(fileName, ".") || strings.Contains(fileName, "~") ||
		   strings.Contains(path, ".tmp.") {
			fmt.Printf("  Skipping temporary file: %s\n", path)
			return nil
		}

		fmt.Printf("  Processing page: %s\n", path)
		return pageProcessor.ProcessPage(path)
	})

	if err != nil {
		return fmt.Errorf("page processing error: %w", err)
	}

	fmt.Printf("  Processing assets from %s to %s\n", cfg.AssetsDir, cfg.PublicDir)
	err = processor.ProcessAssets(cfg.AssetsDir, cfg.PublicDir, false)
	if err != nil {
		log.Printf("Asset processing warning: %v", err)
	}

	return nil
}

func writeErrorToIndex(publicDir, errorMsg string) {
	errorHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Build Error</title>
    <style>
        body { font-family: monospace; margin: 2rem; background: #1a1a1a; color: #ff6b6b; }
        .error { background: #2d1b1b; padding: 1rem; border-radius: 4px; border-left: 4px solid #ff6b6b; }
        pre { white-space: pre-wrap; word-wrap: break-word; }
        .refresh { color: #4ecdc4; margin-top: 1rem; }
    </style>
</head>
<body>
    <h1>🔥 Build Error</h1>
    <div class="error">
        <pre>%s</pre>
    </div>
    <div class="refresh">The page will automatically refresh when you fix the error.</div>
    <script>
        setTimeout(() => location.reload(), 2000);
    </script>
</body>
</html>`, errorMsg)

	indexPath := filepath.Join(publicDir, "index.html")
	os.WriteFile(indexPath, []byte(errorHTML), 0644)
}

func watchAndRebuild(targetDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	cfg := config.NewConfig(targetDir)
	
	watchDirs := []string{cfg.PagesDir, cfg.TemplateDir, cfg.AssetsDir}
	for _, dir := range watchDirs {
		if _, err := os.Stat(dir); err == nil {
			filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					watcher.Add(path)
				}
				return nil
			})
		}
	}

	var lastBuild time.Time
	debounceInterval := 100 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			
			if time.Since(lastBuild) < debounceInterval {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				fileName := filepath.Base(event.Name)
				if strings.Contains(fileName, ".tmp") || strings.HasSuffix(fileName, ".swp") || 
				   strings.HasPrefix(fileName, ".") || strings.Contains(fileName, "~") ||
				   strings.Contains(event.Name, ".tmp.") {
					continue
				}

				fmt.Printf("File changed: %s - rebuilding...\n", event.Name)
				lastBuild = time.Now()

				if err := buildSite(targetDir); err != nil {
					fmt.Printf("Build error: %v\n", err)
					writeErrorToIndex(cfg.PublicDir, err.Error())
				} else {
					fmt.Println("Rebuild completed successfully!")
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve [directory]",
	Short: "Serve the built site locally with live rebuild",
	Long: `Serve the built static site on a local development server with live rebuild.

This command serves files from the public/ directory and automatically
rebuilds the site when source files change. If there are build errors,
an error page will be displayed in the browser.

The server watches for changes in:
- pages/     - Source files (.md, .html, .tmpl)
- templates/ - Template files
- assets/    - Static assets

If the public/ directory doesn't exist, an initial build will be attempted.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		publicDir := filepath.Join(targetDir, "public")
		
		// Try to build the site if public directory doesn't exist
		if _, err := os.Stat(publicDir); os.IsNotExist(err) {
			fmt.Println("Public directory not found, building site...")
			if err := buildSite(targetDir); err != nil {
				return fmt.Errorf("initial build failed: %w", err)
			}
		}

		addr := serveHost + ":" + servePort
		
		fmt.Printf("Serving site from: %s\n", publicDir)
		fmt.Printf("Server running at: http://%s\n", addr)
		fmt.Println("Watching for file changes...")
		fmt.Println("Press Ctrl+C to stop")

		// Start file watcher in a goroutine
		go watchAndRebuild(targetDir)

		fs := http.FileServer(http.Dir(publicDir))
		http.Handle("/", fs)

		if err := http.ListenAndServe(addr, nil); err != nil {
			return fmt.Errorf("server error: %w", err)
		}

		return nil
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8080", "port to serve on")
	serveCmd.Flags().StringVar(&serveHost, "host", "localhost", "host to serve on")
}