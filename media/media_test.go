package media

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yibaiba/wecom"
)

func TestUploadBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/media/upload":
			if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				t.Errorf("content-type %s", r.Header.Get("Content-Type"))
			}
			_, _ = io.Copy(io.Discard, r.Body)
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "type": "file", "media_id": "mid"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	api := New(&wecom.Client{CorpID: "ww", Secret: "s", HTTP: wecom.HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL})
	got, err := api.UploadBytes(context.Background(), "file", "a.txt", []byte("hello"))
	if err != nil || got.MediaID != "mid" {
		t.Fatalf("%+v %v", got, err)
	}
}
