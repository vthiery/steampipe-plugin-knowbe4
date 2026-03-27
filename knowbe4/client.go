package knowbe4

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// ErrNotFound is returned when the KnowBe4 API responds with a 404.
var ErrNotFound = errors.New("not found")

// ErrRateLimited is returned when the KnowBe4 API responds with a 429.
var ErrRateLimited = errors.New("rate limited")

// regionBaseURL maps an API region code to its base URL.
var regionBaseURL = map[string]string{
	"us": "https://us.api.knowbe4.com",
	"eu": "https://eu.api.knowbe4.com",
	"ca": "https://ca.api.knowbe4.com",
	"uk": "https://uk.api.knowbe4.com",
	"de": "https://de.api.knowbe4.com",
}

// Client wraps the KnowBe4 REST Reporting API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// newClient creates a new Client using the given API key and region.
func newClient(apiKey, region string) *Client {
	base, ok := regionBaseURL[region]
	if !ok {
		base = regionBaseURL["us"]
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: base,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getClient retrieves a configured Client from the plugin connection config.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.APIKey == nil {
		return nil, fmt.Errorf("api_key must be configured for the knowbe4 plugin")
	}
	region := "us"
	if cfg.APIRegion != nil && *cfg.APIRegion != "" {
		region = *cfg.APIRegion
	}
	return newClient(*cfg.APIKey, region), nil
}

// get performs an authenticated GET request, unmarshals the response body into result,
// and returns the next cursor value from the X-Next-Cursor response header. An empty string
// means there are no further pages.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// fallthrough to unmarshal
	case http.StatusNotFound:
		return "", ErrNotFound
	case http.StatusTooManyRequests:
		return "", ErrRateLimited
	default:
		return "", fmt.Errorf("KnowBe4 API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}
	return resp.Header.Get("X-Next-Cursor"), nil
}
