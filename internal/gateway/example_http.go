package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ExampleHTTP struct {
	baseURL string
	client  *http.Client
}

func NewExampleHTTP(baseURL string) *ExampleHTTP {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://example.com"
	}

	return &ExampleHTTP{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (g *ExampleHTTP) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+"/", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
