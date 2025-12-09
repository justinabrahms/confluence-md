package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/justinabrahms/confluence-md/internal/confluence"
	md "github.com/JohannesKaufmann/html-to-markdown"
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

// preprocessConfluenceTasks converts Confluence ac:task elements to HTML checkboxes
// before the markdown converter processes them.
func preprocessConfluenceTasks(html string) string {
	// First, replace task-list wrappers with ul tags (before processing tasks)
	result := regexp.MustCompile(`<ac:task-list[^>]*>`).ReplaceAllString(html, "<ul>")
	result = strings.ReplaceAll(result, "</ac:task-list>", "</ul>")

	// Pattern to match individual tasks
	// <ac:task>...<ac:task-status>complete|incomplete</ac:task-status>...<ac:task-body>...content...</ac:task-body>...</ac:task>
	taskPattern := regexp.MustCompile(`(?s)<ac:task>.*?<ac:task-status>(\w+)</ac:task-status>.*?<ac:task-body[^>]*>(.*?)</ac:task-body>.*?</ac:task>`)

	result = taskPattern.ReplaceAllStringFunc(result, func(match string) string {
		submatches := taskPattern.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		status := submatches[1]
		body := submatches[2]

		// Clean up the body - remove span wrappers and inline comment markers
		body = regexp.MustCompile(`<span[^>]*>`).ReplaceAllString(body, "")
		body = strings.ReplaceAll(body, "</span>", "")
		body = regexp.MustCompile(`<ac:inline-comment-marker[^>]*>`).ReplaceAllString(body, "")
		body = strings.ReplaceAll(body, "</ac:inline-comment-marker>", "")
		body = strings.TrimSpace(body)

		checkbox := "[ ]"
		if status == "complete" {
			checkbox = "[x]"
		}

		return fmt.Sprintf("<li>%s %s</li>", checkbox, body)
	})

	return result
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

	// Preprocess Confluence-specific elements
	htmlContent = preprocessConfluenceTasks(htmlContent)

	markdown, err := c.converter.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("converting HTML to markdown: %w", err)
	}

	// Fix escaped checkbox brackets from markdown conversion
	markdown = strings.ReplaceAll(markdown, `\[x\]`, `[x]`)
	markdown = strings.ReplaceAll(markdown, `\[ \]`, `[ ]`)

	output.WriteString(markdown)

	return output.String(), nil
}
