package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/mx-psi/grater/internal/config"
)

type ModuleDep string

type Client struct {
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	return &Client{http.DefaultClient}, nil
}

func url(mod config.ModuleConfig) string {
	return fmt.Sprintf("https://pkg.go.dev/%s@%s?tab=importedby", mod.Path, mod.BaseVersion)
}

// Dependents gets the dependents of a Go module according to pkg.go.dev
func (c *Client) Dependents(ctx context.Context, mod config.ModuleConfig) ([]ModuleDep, error) {
	url := url(mod)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %w", url, err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2xx status code while fetching %q: %d", url, resp.StatusCode)
	}

	return dependentsFromReader(resp.Body)
}

func dependentsFromReader(reader io.Reader) ([]ModuleDep, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	rawPaths := doc.Find(".ImportedBy-detailsIndent").Map(
		func(_ int, s *goquery.Selection) string {
			return s.Text()
		},
	)

	var modulePaths []ModuleDep
	for _, path := range rawPaths {
		modulePaths = append(modulePaths, ModuleDep(path))
	}

	return modulePaths, nil
}
