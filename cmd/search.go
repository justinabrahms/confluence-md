package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/confluence-md/internal/confluence"
	"github.com/yourusername/confluence-md/internal/config"
	"github.com/yourusername/confluence-md/internal/markdown"
)

var (
	spaceKey   string
	limit      int
	lucky      bool
	resultIndex int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for Confluence pages",
	Long:  `Search for Confluence pages by query string and optionally fetch the content.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading configuration: %w", err)
		}

		// Create client
		client := confluence.NewClient(cfg.ConfluenceURL, cfg.Email, cfg.APIToken, Debug)

		if Debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Config: URL=%s, Email=%s\n", cfg.ConfluenceURL, cfg.Email)
			fmt.Fprintf(os.Stderr, "[DEBUG] Query: %s, Space: %s, Limit: %d\n", query, spaceKey, limit)
		}

		// Search
		results, err := client.Search(query, spaceKey, limit)
		if err != nil {
			return fmt.Errorf("searching: %w", err)
		}

		if results.Size == 0 {
			fmt.Println("No results found")
			os.Exit(4)
		}

		// If --lucky or --index is specified, fetch the content
		if lucky || resultIndex > 0 {
			fetchIndex := 0
			if resultIndex > 0 {
				fetchIndex = resultIndex - 1 // Convert to 0-based
			}

			if fetchIndex >= len(results.Results) {
				return fmt.Errorf("index %d out of range (found %d results)", resultIndex, len(results.Results))
			}

			result := results.Results[fetchIndex]
			if Debug {
				fmt.Fprintf(os.Stderr, "[DEBUG] Selected result: Title=%s, ID=%s, Type=%s\n",
					result.Title, result.ID, result.Type)
			}
			page, err := client.GetPageByID(result.ID)
			if err != nil {
				return fmt.Errorf("fetching page: %w", err)
			}

			converter := markdown.NewConverter()
			md, err := converter.PageToMarkdown(page, includeMetadata)
			if err != nil {
				return fmt.Errorf("converting to markdown: %w", err)
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(md), 0644); err != nil {
					return fmt.Errorf("writing to file: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Written to %s\n", outputFile)
			} else {
				fmt.Print(md)
			}

			return nil
		}

		// Display search results
		fmt.Printf("Found %d results:\n\n", results.Size)

		for i, result := range results.Results {
			pageURL := cfg.ConfluenceURL + result.Links.WebUI

			fmt.Printf("[%d] %s\n", i+1, result.Title)
			fmt.Printf("    Space: %s | Updated: %s\n",
				result.Space.Key,
				result.LastModified.Format("2006-01-02"))
			fmt.Printf("    URL: %s\n\n", pageURL)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&spaceKey, "space", "", "Limit search to specific space")
	searchCmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results (1-50)")
	searchCmd.Flags().BoolVar(&lucky, "lucky", false, "Fetch the first search result")
	searchCmd.Flags().IntVar(&resultIndex, "index", 0, "Fetch a specific search result by index (1-based)")
	searchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
	searchCmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Include page metadata in output")
}
