package knowbe4

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"golang.org/x/time/rate"
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

// KnowBe4 caps Reporting API usage at 4 req/s and 50 req/min burst. The
// tightest sustained ceiling is 50/min, so the limiter paces at that rate
// while still allowing short bursts up to the documented burst budget.
const (
	rateLimitPerMinute = 50
	rateLimitBurst     = 4
	maxRetryAttempts   = 5
	maxBackoff         = 60 * time.Second
)

// Client wraps the KnowBe4 REST Reporting API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	limiter    *rate.Limiter
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
		limiter: rate.NewLimiter(rate.Every(time.Minute/rateLimitPerMinute), rateLimitBurst),
	}
}

// getClient retrieves a configured Client from the plugin connection config.
// The client is cached on the connection so its rate limiter is shared across
// every concurrent hydrate for the same connection.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.APIKey == nil {
		return nil, fmt.Errorf("api_key must be configured for the knowbe4 plugin")
	}
	region := "us"
	if cfg.APIRegion != nil && *cfg.APIRegion != "" {
		region = *cfg.APIRegion
	}

	cacheKey := "knowbe4_client:" + region
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		if client, ok := cached.(*Client); ok {
			return client, nil
		}
	}
	client := newClient(*cfg.APIKey, region)
	d.ConnectionManager.Cache.Set(cacheKey, client)
	return client, nil
}

// get performs an authenticated GET request and unmarshals the response body
// into result. Requests are paced through the client's rate limiter, and 429
// responses are retried (honoring Retry-After when present, otherwise using
// exponential backoff) up to maxRetryAttempts times.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	var lastErr error
	for attempt := 0; attempt < maxRetryAttempts; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}

		retryAfter, err := c.do(ctx, path, params, result)
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrRateLimited) {
			return err
		}
		lastErr = err

		delay := retryAfter
		if delay <= 0 {
			delay = time.Duration(1<<attempt) * time.Second
		}
		if delay > maxBackoff {
			delay = maxBackoff
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return lastErr
}

// do issues a single GET. On 429 it returns ErrRateLimited along with any
// Retry-After hint parsed from the response.
func (c *Client) do(ctx context.Context, path string, params map[string]string, result interface{}) (time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
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
		return 0, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// fallthrough to unmarshal
	case http.StatusNotFound:
		return 0, ErrNotFound
	case http.StatusTooManyRequests:
		return parseRetryAfter(resp.Header.Get("Retry-After")), ErrRateLimited
	default:
		return 0, fmt.Errorf("KnowBe4 API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return 0, fmt.Errorf("unmarshalling response: %w", err)
	}
	return 0, nil
}

// parseRetryAfter parses an HTTP Retry-After header value, which may be either
// an integer number of seconds or an HTTP-date. Returns 0 when the value is
// missing or unparseable so the caller falls back to exponential backoff.
func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// cursorResponse is the envelope KnowBe4 returns when cursor pagination is enabled
// (i.e. when the request includes `cursor=true`).
type cursorResponse[T any] struct {
	Data     []T `json:"data"`
	Metadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"_metadata"`
}

// listPaged performs a cursor-paginated GET, returning the page items and the next
// cursor token. An empty next cursor signals there are no further pages.
func listPaged[T any](ctx context.Context, c *Client, path string, params map[string]string) ([]T, string, error) {
	var resp cursorResponse[T]
	if err := c.get(ctx, path, params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Data, resp.Metadata.NextCursor, nil
}
