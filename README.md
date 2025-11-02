# Confluence Markdown Fetcher

A command-line tool for retrieving Confluence pages as Markdown, designed for both human readers and LLM consumption.

## Features

- **Fetch by URL**: Download a Confluence page directly using its URL
- **Search and fetch**: Find pages by name/query and retrieve their content
- **Markdown output**: Pages are converted to clean, readable Markdown
- **Machine-readable format**: Output structured for easy parsing by LLMs and automation tools

## Installation

Download the latest binary release for your platform from the [releases page](https://github.com/yourusername/confluence-md/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/yourusername/confluence-md/releases/latest/download/confluence-md-darwin-arm64 -o confluence-md
chmod +x confluence-md
sudo mv confluence-md /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yourusername/confluence-md/releases/latest/download/confluence-md-darwin-amd64 -o confluence-md
chmod +x confluence-md
sudo mv confluence-md /usr/local/bin/

# Linux
curl -L https://github.com/yourusername/confluence-md/releases/latest/download/confluence-md-linux-amd64 -o confluence-md
chmod +x confluence-md
sudo mv confluence-md /usr/local/bin/
```

## Configuration

Configure your Confluence credentials using either a config file or environment variables.

### Config File (Recommended)

Create a config file at `~/.config/confluence-md/config.yaml`:

```yaml
confluence_url: https://your-domain.atlassian.net/wiki
email: your-email@example.com
api_token: your-api-token
```

You can also set `XDG_CONFIG_HOME` to use a different config directory:
```bash
export XDG_CONFIG_HOME=/custom/path
# Config file will be at /custom/path/confluence-md/config.yaml
```

### Environment Variables

Environment variables override values in the config file:

```bash
export CONFLUENCE_URL="https://your-domain.atlassian.net/wiki"
export CONFLUENCE_EMAIL="your-email@example.com"
export CONFLUENCE_API_TOKEN="your-api-token"
```

### Getting an API Token

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a label and copy the token
4. Add it to your config file or environment

## Usage

### Fetch a page by URL

```bash
confluence-md fetch https://your-domain.atlassian.net/wiki/spaces/TEAM/pages/123456/Page+Title
```

Output is written to stdout as Markdown.

### Search for pages

```bash
# Search and list results
confluence-md search "project documentation"
```

Example output:
```
Found 5 results:

[1] Project Documentation Overview
    Space: ENG | Updated: 2025-01-15
    URL: https://company.atlassian.net/wiki/spaces/ENG/pages/123456/Project+Documentation+Overview

[2] API Documentation Standards
    Space: ENG | Updated: 2025-01-10
    URL: https://company.atlassian.net/wiki/spaces/ENG/pages/123789/API+Documentation+Standards

[3] Documentation Review Process
    Space: TEAM | Updated: 2024-12-20
    URL: https://company.atlassian.net/wiki/spaces/TEAM/pages/456123/Documentation+Review+Process

[4] Legacy Project Docs
    Space: ARCHIVE | Updated: 2024-06-15
    URL: https://company.atlassian.net/wiki/spaces/ARCHIVE/pages/789456/Legacy+Project+Docs

[5] Project Roadmap Documentation
    Space: PRODUCT | Updated: 2025-01-05
    URL: https://company.atlassian.net/wiki/spaces/PRODUCT/pages/321654/Project+Roadmap+Documentation
```

### Fetch from search results

```bash
# Search and fetch the first result (I'm feeling lucky!)
confluence-md search "project documentation" --lucky

# Search and fetch a specific result by index
confluence-md search "project documentation" --index 2

# Or pipe the search result URL to fetch
confluence-md search "onboarding" | grep "^\[2\]" | confluence-md fetch
```

### Options

- `--output, -o`: Write output to a file instead of stdout
- `--space`: Limit search to a specific Confluence space
- `--limit`: Maximum number of search results to return (default: 10)
- `--include-metadata`: Include page metadata (author, dates, labels) in output
- `--lucky`: Automatically fetch content from the first search result
- `--index`: Which search result to fetch (1-based index)

## Examples

```bash
# Fetch page and save to file
confluence-md fetch https://company.atlassian.net/wiki/spaces/ENG/pages/789/API-Docs -o api-docs.md

# Search in specific space
confluence-md search "onboarding" --space HR

# Search and fetch first result with metadata (I'm feeling lucky!)
confluence-md search "architecture decision" --lucky --include-metadata

# Fetch the third search result
confluence-md search "API guidelines" --index 3

# Pipe to other tools
confluence-md search "sprint planning" --lucky | grep TODO
```

## Output Format

### Search results
Numbered list with page title, space, last updated date, and full URL for easy reference.

### Markdown content
- Page title as H1
- Optional metadata block (when `--include-metadata` is used)
- Page content converted to Markdown
- Links preserved and converted to Markdown format
- Code blocks, tables, and formatting maintained

## Exit Codes

- `0`: Success
- `1`: General error (invalid arguments, network errors)
- `2`: Authentication failure
- `3`: Page not found
- `4`: No search results found
