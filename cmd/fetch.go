package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/justinabrahms/confluence-md/internal/confluence"
	"github.com/justinabrahms/confluence-md/internal/config"
	"github.com/justinabrahms/confluence-md/internal/markdown"
)

var (
	outputFile      string
	includeMetadata bool
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [url]",
	Short: "Fetch a Confluence page by URL",
	Long:  `Fetch a Confluence page by its URL and convert it to Markdown.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pageURL := args[0]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading configuration: %w", err)
		}

		// Create client
		client := confluence.NewClient(cfg.ConfluenceURL, cfg.Email, cfg.APIToken, Debug)

		if Debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Config: URL=%s, Email=%s\n", cfg.ConfluenceURL, cfg.Email)
			fmt.Fprintf(os.Stderr, "[DEBUG] Fetching URL: %s\n", pageURL)
		}

		// Fetch page
		page, err := client.GetPageByURL(pageURL)
		if err != nil {
			return fmt.Errorf("fetching page: %w", err)
		}

		// Convert to markdown
		converter := markdown.NewConverter()
		md, err := converter.PageToMarkdown(page, includeMetadata)
		if err != nil {
			return fmt.Errorf("converting to markdown: %w", err)
		}

		// Output
		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(md), 0644); err != nil {
				return fmt.Errorf("writing to file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Written to %s\n", outputFile)
		} else {
			fmt.Print(md)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
	fetchCmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Include page metadata in output")
}
