package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// UploadedMedia is the result of media/upload.
type UploadedMedia struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type multiparter interface {
	PostMultipart(ctx context.Context, rawURL, field, filename string, r io.Reader) ([]byte, error)
}

// UploadMedia uploads a temporary media file. typ is image, voice, video, or file.
func (c *Client) UploadMedia(ctx context.Context, typ, filename string, r io.Reader) (UploadedMedia, error) {
	tok, err := c.Token(ctx)
	if err != nil {
		return UploadedMedia{}, err
	}
	q := url.Values{}
	q.Set("access_token", tok)
	q.Set("type", typ)
	raw := apiBase + "/cgi-bin/media/upload?" + q.Encode()
	body, err := c.multipart(ctx, raw, "media", filename, r)
	if err != nil {
		return UploadedMedia{}, err
	}
	var out struct {
		apiMeta
		UploadedMedia
	}
	if err := decodeAPI(body, &out); err != nil {
		return UploadedMedia{}, err
	}
	return out.UploadedMedia, nil
}

// UploadImage uploads a permanent image (media/uploadimg) and returns its URL.
func (c *Client) UploadImage(ctx context.Context, filename string, r io.Reader) (string, error) {
	tok, err := c.Token(ctx)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	q.Set("access_token", tok)
	raw := apiBase + "/cgi-bin/media/uploadimg?" + q.Encode()
	body, err := c.multipart(ctx, raw, "media", filename, r)
	if err != nil {
		return "", err
	}
	var out struct {
		apiMeta
		URL string `json:"url"`
	}
	if err := decodeAPI(body, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}

// GetMedia downloads temporary media. On API errors it returns Error.
func (c *Client) GetMedia(ctx context.Context, mediaID string) ([]byte, error) {
	tok, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("access_token", tok)
	q.Set("media_id", mediaID)
	body, err := c.doer().Get(ctx, apiBase+"/cgi-bin/media/get?"+q.Encode())
	if err != nil {
		return nil, err
	}
	if len(body) > 0 && body[0] == '{' {
		var meta apiMeta
		if err := json.Unmarshal(body, &meta); err == nil && meta.ErrCode != 0 {
			return nil, Error{Code: meta.ErrCode, Msg: meta.ErrMsg}
		}
	}
	return body, nil
}

func (c *Client) multipart(ctx context.Context, rawURL, field, filename string, r io.Reader) ([]byte, error) {
	m, ok := c.doer().(multiparter)
	if !ok {
		return nil, fmt.Errorf("wecom http client does not support multipart upload")
	}
	return m.PostMultipart(ctx, rawURL, field, filename, r)
}

// UploadMediaBytes is UploadMedia from a byte slice.
func (c *Client) UploadMediaBytes(ctx context.Context, typ, filename string, data []byte) (UploadedMedia, error) {
	return c.UploadMedia(ctx, typ, filename, bytes.NewReader(data))
}
