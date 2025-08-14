package commands

import (
	"fmt"

	"github.com/ahoglund/go-static/pkg/scaffold"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new static site",
	Long: `Initialize a new static site with the basic directory structure 
and example files.

Creates:
- pages/     - Directory for markdown and HTML files
- templates/ - Directory for Go template files  
- assets/    - Directory for static assets
- Example templates and a sample page`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		scaffolder := scaffold.NewScaffolder(verbose)

		if err := scaffolder.CreateSite(targetDir); err != nil {
			return fmt.Errorf("failed to create site: %w", err)
		}

		if err := scaffolder.CreateDirectories(targetDir); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}

		fmt.Printf("Site initialized successfully in: %s\n", targetDir)
		fmt.Println("\nNext steps:")
		if targetDir != "." {
			fmt.Printf("  cd %s\n", targetDir)
		}
		fmt.Println("  go-static build")
		fmt.Println("  go-static serve")

		return nil
	},
}