package knowbe4

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

type item struct {
	ID int `json:"id"`
}

func testClient(url string) *Client {
	return &Client{
		apiKey:     "test",
		baseURL:    url,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		limiter:    rate.NewLimiter(rate.Inf, 1),
	}
}

func TestListPagedReadsNextCursor(t *testing.T) {
	cases := map[string]struct {
		body string
		want string
	}{
		"pagination.nextCursor": {`{"data":[{"id":1}],"pagination":{"perPage":500,"nextCursor":"abc"}}`, "abc"},
		"_metadata.next_cursor": {`{"data":[{"id":1}],"_metadata":{"next_cursor":"def"}}`, "def"},
		"last page":             {`{"data":[{"id":1}],"pagination":{"perPage":500}}`, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprint(w, tc.body)
			}))
			defer srv.Close()

			data, next, err := listPaged[item](context.Background(), testClient(srv.URL), "/v1/users", nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(data) != 1 || next != tc.want {
				t.Fatalf("got %d items, cursor %q; want 1 item, cursor %q", len(data), next, tc.want)
			}
		})
	}
}

// Follows nextCursor across pages the way the table list functions do.
func TestListPagedFollowsAllPages(t *testing.T) {
	pages := map[string]string{
		"true": `{"data":[{"id":1},{"id":2}],"pagination":{"perPage":2,"nextCursor":"p2"}}`,
		"p2":   `{"data":[{"id":3},{"id":4}],"pagination":{"perPage":2,"nextCursor":"p3"}}`,
		"p3":   `{"data":[{"id":5}],"pagination":{"perPage":2}}`,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := pages[r.URL.Query().Get("cursor")]
		if !ok {
			http.Error(w, "unknown cursor", http.StatusBadRequest)
			return
		}
		_, _ = fmt.Fprint(w, body)
	}))
	defer srv.Close()

	client := testClient(srv.URL)
	params := map[string]string{"cursor": "true", "per_page": "2"}
	var all []item
	for {
		data, next, err := listPaged[item](context.Background(), client, "/v1/users", params)
		if err != nil {
			t.Fatal(err)
		}
		all = append(all, data...)
		if next == "" || len(data) == 0 {
			break
		}
		params["cursor"] = next
	}
	if len(all) != 5 {
		t.Fatalf("got %d items across pages, want 5", len(all))
	}
}
