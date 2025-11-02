# Confluence Document Reader

Read and search Confluence pages as markdown using the confluence-md CLI tool.

## What this skill does

This skill allows Claude to:
- Search for Confluence pages by query
- Fetch specific Confluence pages by URL
- Filter searches by space or creator
- Return content as clean markdown for easy reading and analysis

## When to use this skill

Use this skill when the user asks to:
- Read a Confluence page
- Search for Confluence documentation
- Find pages about a specific topic
- Get content from a Confluence URL
- Review their own Confluence pages (with --mine flag)

## Usage

The skill uses the `confluence-md` command-line tool installed on the system.

### Search for pages

```bash
confluence-md search "query" [flags]
```

Common flags:
- `--lucky` - Fetch the first result immediately
- `--mine` - Only search pages you created
- `--space SPACE` - Limit search to specific space
- `--limit N` - Maximum number of results (default 10)
- `--index N` - Fetch specific result by number

### Fetch a specific page

```bash
confluence-md fetch <url>
```

## Examples

```bash
# Search and get first result about API documentation
confluence-md search "API documentation" --lucky

# Find your own proposals
confluence-md search "proposal" --mine --limit 5

# Search in specific space
confluence-md search "onboarding" --space HR --lucky

# Get specific page
confluence-md fetch https://company.atlassian.net/wiki/spaces/ENG/pages/123456/Page+Title
```

## Output

All commands return markdown-formatted content that can be directly analyzed and discussed.

## Requirements

- The `confluence-md` binary must be in your PATH
- Configuration must be set up in `~/.config/confluence-md/config.yaml` or via environment variables
- Valid Confluence credentials (email and API token)

## Tips

- Use `--lucky` when you want the most relevant result quickly
- Use `--mine` to find pages you authored
- The tool converts HTML to clean markdown automatically
- You can pipe output to other tools for further processing
