package wecom

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

// HTTPDoer performs GET and POST requests against the WeCom API.
type HTTPDoer struct {
	Client *http.Client
}

func (h HTTPDoer) client() *http.Client {
	if h.Client != nil {
		return h.Client
	}
	return &http.Client{Timeout: 10 * time.Second}
}

// Get fetches a URL.
func (h HTTPDoer) Get(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	return h.do(req)
}

// Post sends JSON to a URL.
func (h HTTPDoer) Post(ctx context.Context, rawURL string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.do(req)
}

func (h HTTPDoer) do(req *http.Request) ([]byte, error) {
	resp, err := h.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}
