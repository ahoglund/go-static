package commands

import (
	"fmt"

	"github.com/ahoglund/go-static/pkg/scaffold"
	"github.com/spf13/cobra"
)

var githubPages bool

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new static site",
	Long: `Initialize a new static site with the basic directory structure 
and example files.

Creates:
- pages/     - Directory for markdown and HTML files
- templates/ - Directory for Go template files  
- assets/    - Directory for static assets
- Example templates and sample pages

With --github-pages flag:
- .gitignore - Ignores public/ directory
- .github/workflows/deploy.yml - Auto-deploy to GitHub Pages
- README.md - Documentation for GitHub repository`,
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

		if githubPages {
			if err := scaffolder.CreateGitHubPages(targetDir); err != nil {
				return fmt.Errorf("failed to create GitHub Pages files: %w", err)
			}
		}

		fmt.Printf("Site initialized successfully in: %s\n", targetDir)
		
		if githubPages {
			fmt.Println("\n🚀 GitHub Pages setup complete!")
			fmt.Println("Next steps:")
			if targetDir != "." {
				fmt.Printf("  cd %s\n", targetDir)
			}
			fmt.Println("  git init")
			fmt.Println("  git add .")
			fmt.Println("  git commit -m \"Initial commit\"")
			fmt.Println("  git branch -M main")
			fmt.Println("  git remote add origin <your-repo-url>")
			fmt.Println("  git push -u origin main")
			fmt.Println("\nThen enable GitHub Pages in your repository settings!")
		} else {
			fmt.Println("\nNext steps:")
			if targetDir != "." {
				fmt.Printf("  cd %s\n", targetDir)
			}
			fmt.Println("  go-static build")
			fmt.Println("  go-static serve")
		}

		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&githubPages, "github-pages", false, "setup GitHub Pages deployment with .gitignore and workflow")
}