package commands

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	servePort string
	serveHost string
)

var serveCmd = &cobra.Command{
	Use:   "serve [directory]",
	Short: "Serve the built site locally",
	Long: `Serve the built static site on a local development server.

This command serves files from the public/ directory. You should run 
'go-static build' first to generate the static files.

For development, the server will serve files directly without 
live reload (coming in future versions).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		publicDir := filepath.Join(targetDir, "public")
		
		if _, err := os.Stat(publicDir); os.IsNotExist(err) {
			return fmt.Errorf("public directory not found: %s\nRun 'go-static build' first", publicDir)
		}

		addr := serveHost + ":" + servePort
		
		fmt.Printf("Serving site from: %s\n", publicDir)
		fmt.Printf("Server running at: http://%s\n", addr)
		fmt.Println("Press Ctrl+C to stop")

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