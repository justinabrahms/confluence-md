package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Email      string
	APIToken   string
	HTTPClient *http.Client
	Debug      bool
	logger     *log.Logger
}

type Page struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Title   string    `json:"title"`
	Body    Body      `json:"body"`
	Version Version   `json:"version"`
	History History   `json:"history"`
	Space   Space     `json:"_expandable"`
	Links   Links     `json:"_links"`
}

type Body struct {
	Storage Storage `json:"storage"`
	View    Storage `json:"view"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Version struct {
	Number  int       `json:"number"`
	When    time.Time `json:"when"`
	Message string    `json:"message"`
	By      User      `json:"by"`
}

type User struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type History struct {
	CreatedDate time.Time `json:"createdDate"`
	CreatedBy   User      `json:"createdBy"`
}

type Space struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Links struct {
	WebUI string `json:"webui"`
	Self  string `json:"self"`
}

type SearchResult struct {
	Results []SearchResultItem `json:"results"`
	Start   int                `json:"start"`
	Limit   int                `json:"limit"`
	Size    int                `json:"size"`
}

type SearchResultItem struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	Space         Space     `json:"space"`
	History       History   `json:"history"`
	Version       Version   `json:"version"`
	Links         Links     `json:"_links"`
	Excerpt       string    `json:"excerpt"`
	URL           string    `json:"url"`
	LastModified  time.Time `json:"lastModified"`
	FriendlyDate  string    `json:"friendlyLastModified"`
}

func NewClient(baseURL, email, apiToken string, debug bool) *Client {
	logger := log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
	return &Client{
		BaseURL:  strings.TrimSuffix(baseURL, "/"),
		Email:    email,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Debug:  debug,
		logger: logger,
	}
}

func (c *Client) debugf(format string, args ...interface{}) {
	if c.Debug {
		c.logger.Printf(format, args...)
	}
}

func (c *Client) doRequest(method, path string) (*http.Response, error) {
	fullURL := c.BaseURL + path
	c.debugf("Request: %s %s", method, fullURL)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.Email, c.APIToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	c.debugf("Response: HTTP %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		c.debugf("Error response body: %s", string(body))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) GetPageByID(pageID string) (*Page, error) {
	c.debugf("Fetching page by ID: %s", pageID)
	path := fmt.Sprintf("/rest/api/content/%s?expand=body.storage,body.view,version,history,space", pageID)

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	c.debugf("Successfully fetched page: %s (ID: %s)", page.Title, page.ID)
	return &page, nil
}

func (c *Client) GetPageByURL(pageURL string) (*Page, error) {
	pageID, err := extractPageIDFromURL(pageURL)
	if err != nil {
		return nil, err
	}
	return c.GetPageByID(pageID)
}

func (c *Client) Search(query string, spaceKey string, limit int) (*SearchResult, error) {
	params := url.Values{}

	cql := query
	if spaceKey != "" {
		cql = fmt.Sprintf("type=page AND space=%s AND text~\"%s\"", spaceKey, query)
	} else {
		cql = fmt.Sprintf("type=page AND text~\"%s\"", query)
	}

	c.debugf("Search CQL: %s", cql)

	params.Set("cql", cql)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("expand", "space,version,history,lastModified")

	path := "/rest/api/content/search?" + params.Encode()

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if c.Debug {
		c.debugf("Raw search response (first 500 chars): %s", string(body[:min(500, len(body))]))
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding search results: %w", err)
	}

	c.debugf("Search returned %d results", result.Size)
	if c.Debug && len(result.Results) > 0 {
		c.debugf("First result: Title=%s, ID=%s", result.Results[0].Title, result.Results[0].ID)
	}
	return &result, nil
}

func extractPageIDFromURL(pageURL string) (string, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	parts := strings.Split(u.Path, "/pages/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid Confluence URL format: missing /pages/ segment")
	}

	idPart := strings.Split(parts[1], "/")[0]
	if idPart == "" {
		return "", fmt.Errorf("could not extract page ID from URL")
	}

	return idPart, nil
}
