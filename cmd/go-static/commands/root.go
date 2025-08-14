package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"  // Can be set via build ldflags
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "go-static",
	Short: "A fast static site generator written in Go",
	Long: `go-static is a static site generator that converts markdown files 
with frontmatter into HTML pages using Go templates.

It follows simple conventions for directory structure and provides
a fast, efficient way to build static websites.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}