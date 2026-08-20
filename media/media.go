package media

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"

	"github.com/yibaiba/wecom"
)

// API uploads and downloads temporary media.
type API struct {
	App *wecom.Client
}

// New wraps a shared WeCom client.
func New(app *wecom.Client) *API {
	return &API{App: app}
}

// Uploaded is the result of media/upload.
type Uploaded struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// Upload uploads a temporary media file. typ is image, voice, video, or file.
func (a *API) Upload(ctx context.Context, typ, filename string, r io.Reader) (Uploaded, error) {
	q := url.Values{}
	q.Set("type", typ)
	body, err := a.App.PostFormFile(ctx, "/cgi-bin/media/upload", q, "media", filename, r)
	if err != nil {
		return Uploaded{}, err
	}
	var out Uploaded
	if err := wecom.Decode(body, &out); err != nil {
		return Uploaded{}, err
	}
	return out, nil
}

// UploadImage uploads a permanent image (media/uploadimg) and returns its URL.
func (a *API) UploadImage(ctx context.Context, filename string, r io.Reader) (string, error) {
	body, err := a.App.PostFormFile(ctx, "/cgi-bin/media/uploadimg", nil, "media", filename, r)
	if err != nil {
		return "", err
	}
	var out struct {
		URL string `json:"url"`
	}
	if err := wecom.Decode(body, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}

// Get downloads temporary media. On API errors it returns wecom.Error.
func (a *API) Get(ctx context.Context, mediaID string) ([]byte, error) {
	body, err := a.App.GetRaw(ctx, "/cgi-bin/media/get", url.Values{"media_id": {mediaID}})
	if err != nil {
		return nil, err
	}
	if len(body) > 0 && body[0] == '{' {
		var meta struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(body, &meta); err == nil && meta.ErrCode != 0 {
			return nil, wecom.Error{Code: meta.ErrCode, Msg: meta.ErrMsg}
		}
	}
	return body, nil
}

// UploadBytes is Upload from a byte slice.
func (a *API) UploadBytes(ctx context.Context, typ, filename string, data []byte) (Uploaded, error) {
	return a.Upload(ctx, typ, filename, bytes.NewReader(data))
}
