package markdown

import (
	"fmt"
	"strings"

	"github.com/justinabrahms/confluence-md/internal/confluence"
	"github.com/JohannesKaufmann/html-to-markdown"
)

type Converter struct {
	converter *md.Converter
}

func NewConverter() *Converter {
	converter := md.NewConverter("", true, nil)
	return &Converter{
		converter: converter,
	}
}

func (c *Converter) PageToMarkdown(page *confluence.Page, includeMetadata bool) (string, error) {
	var output strings.Builder

	// Title as H1
	output.WriteString(fmt.Sprintf("# %s\n\n", page.Title))

	// Optional metadata
	if includeMetadata {
		output.WriteString("---\n\n")
		output.WriteString(fmt.Sprintf("**Page ID:** %s\n\n", page.ID))
		output.WriteString(fmt.Sprintf("**Created:** %s by %s\n\n",
			page.History.CreatedDate.Format("2006-01-02"),
			page.History.CreatedBy.DisplayName))
		output.WriteString(fmt.Sprintf("**Last Modified:** %s by %s (v%d)\n\n",
			page.Version.When.Format("2006-01-02"),
			page.Version.By.DisplayName,
			page.Version.Number))
		if page.Version.Message != "" {
			output.WriteString(fmt.Sprintf("**Version Message:** %s\n\n", page.Version.Message))
		}
		output.WriteString("---\n\n")
	}

	// Convert HTML content to Markdown
	htmlContent := page.Body.Storage.Value
	if htmlContent == "" {
		htmlContent = page.Body.View.Value
	}

	markdown, err := c.converter.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("converting HTML to markdown: %w", err)
	}

	output.WriteString(markdown)

	return output.String(), nil
}
