package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &Client{CorpID: "ww", AgentID: 1000036, Secret: "s", HTTP: HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL}
}

func writeToken(w http.ResponseWriter) {
	_ = json.NewEncoder(w).Encode(tokenResp{AccessToken: "tok", ExpiresIn: 7200, ErrCode: 0})
}

func TestClientCachesToken(t *testing.T) {
	tokens := 0
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			tokens++
			writeToken(w)
		case "/cgi-bin/agent/list":
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "agentlist": []any{}})
		default:
			http.NotFound(w, r)
		}
	})
	var out struct {
		AgentList []any `json:"agentlist"`
	}
	if err := cli.GetJSON(context.Background(), "/cgi-bin/agent/list", nil, &out); err != nil {
		t.Fatal(err)
	}
	if err := cli.GetJSON(context.Background(), "/cgi-bin/agent/list", nil, &out); err != nil {
		t.Fatal(err)
	}
	if tokens != 1 {
		t.Fatalf("token fetches %d", tokens)
	}
}

func TestClientAPIError(t *testing.T) {
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			writeToken(w)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 81013, "errmsg": "user out of scope"})
	})
	err := cli.GetJSON(context.Background(), "/cgi-bin/department/list", nil, &struct{}{})
	e, ok := err.(Error)
	if !ok || e.Code != 81013 {
		t.Fatalf("err %v", err)
	}
}

func TestClientTokenRetry(t *testing.T) {
	n := 0
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeToken(w)
		case "/cgi-bin/department/list":
			n++
			if n == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 42001, "errmsg": "expired"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "department": []any{}})
		default:
			http.NotFound(w, r)
		}
	})
	cli.token, cli.tokenExp = "old", time.Now().Add(time.Hour)
	var out struct {
		Department []any `json:"department"`
	}
	if err := cli.GetJSON(context.Background(), "/cgi-bin/department/list", nil, &out); err != nil {
		t.Fatal(err)
	}
}

func TestClientTokenSerializesRefresh(t *testing.T) {
	var tokens atomic.Int32
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/gettoken" {
			http.NotFound(w, r)
			return
		}
		tokens.Add(1)
		time.Sleep(20 * time.Millisecond)
		writeToken(w)
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := cli.Token(context.Background()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := tokens.Load(); n != 1 {
		t.Fatalf("token fetches %d", n)
	}
}

func TestGetRawTokenRetry(t *testing.T) {
	n := 0
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeToken(w)
		case "/cgi-bin/media/get":
			n++
			if n == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 42001, "errmsg": "expired"})
				return
			}
			_, _ = w.Write([]byte("PNG"))
		default:
			http.NotFound(w, r)
		}
	})
	cli.token, cli.tokenExp = "old", time.Now().Add(time.Hour)
	got, err := cli.GetRaw(context.Background(), "/cgi-bin/media/get", url.Values{"media_id": {"m"}})
	if err != nil || string(got) != "PNG" {
		t.Fatalf("%s %v", got, err)
	}
}

func TestPostFormFileTokenRetry(t *testing.T) {
	n := 0
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeToken(w)
		case "/cgi-bin/media/upload":
			n++
			if n == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 42001, "errmsg": "expired"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "media_id": "mid"})
		default:
			http.NotFound(w, r)
		}
	})
	cli.token, cli.tokenExp = "old", time.Now().Add(time.Hour)
	got, err := cli.PostFormFile(context.Background(), "/cgi-bin/media/upload", url.Values{"type": {"file"}}, "media", "a.txt", bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}
	if err := Decode(got, &struct{}{}); err != nil {
		t.Fatal(err)
	}
}
