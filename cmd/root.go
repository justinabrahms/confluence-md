package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Debug bool
)

var rootCmd = &cobra.Command{
	Use:   "confluence-md",
	Short: "Confluence Markdown Fetcher",
	Long: `A command-line tool for retrieving Confluence pages as Markdown.

Fetch pages by URL or search for pages by name and retrieve their content.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug logging")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
