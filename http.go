package wecom

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

var defaultHTTP = &http.Client{Timeout: 10 * time.Second}

// HTTPDoer performs GET and POST requests against the WeCom API.
type HTTPDoer struct {
	Client *http.Client
}

func (h HTTPDoer) client() *http.Client {
	if h.Client != nil {
		return h.Client
	}
	return defaultHTTP
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
	return io.ReadAll(io.LimitReader(resp.Body, 21<<20))
}

// PostMultipart uploads a form file field (media/upload).
func (h HTTPDoer) PostMultipart(ctx context.Context, rawURL, field, filename string, r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, r); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return h.do(req)
}
