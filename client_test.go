package wecom

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	orig := apiBase
	apiBase = srv.URL
	t.Cleanup(func() { apiBase = orig })
	return &Client{CorpID: "ww", AgentID: 1000036, Secret: "s", HTTP: HTTPDoer{Client: srv.Client()}}
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
	if _, err := cli.ListAgents(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := cli.ListAgents(context.Background()); err != nil {
		t.Fatal(err)
	}
	if tokens != 1 {
		t.Fatalf("token fetches %d", tokens)
	}
}

func TestClientGetAgentAndSendText(t *testing.T) {
	var sent map[string]any
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeToken(w)
		case "/cgi-bin/agent/get":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errcode": 0, "agentid": 1000036, "name": "SSO",
				"allow_userinfos": map[string]any{"user": []any{map[string]string{"userid": "zhangsan"}}},
				"allow_partys":    map[string]any{"partyid": []int{1}},
			})
		case "/cgi-bin/message/send":
			_ = json.NewDecoder(r.Body).Decode(&sent)
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "msgid": "m1"})
		default:
			http.NotFound(w, r)
		}
	})
	ag, err := cli.GetAgent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ag.Name != "SSO" || len(ag.AllowUserIDs) != 1 || ag.AllowUserIDs[0] != "zhangsan" {
		t.Fatalf("agent %+v", ag)
	}
	res, err := cli.SendText(context.Background(), "zhangsan", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if res.MsgID != "m1" || sent["msgtype"] != "text" || sent["agentid"].(float64) != 1000036 {
		t.Fatalf("send %+v %+v", res, sent)
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
	_, err := cli.ListDepartments(context.Background(), 0)
	e, ok := err.(Error)
	if !ok || e.Code != 81013 {
		t.Fatalf("err %v", err)
	}
}

func TestClientUploadMedia(t *testing.T) {
	cli := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeToken(w)
		case "/cgi-bin/media/upload":
			if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				t.Errorf("content-type %s", r.Header.Get("Content-Type"))
			}
			_, _ = io.Copy(io.Discard, r.Body)
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "type": "file", "media_id": "mid"})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := cli.UploadMediaBytes(context.Background(), "file", "a.txt", []byte("hello"))
	if err != nil || got.MediaID != "mid" {
		t.Fatalf("upload %+v %v", got, err)
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
	got, err := cli.ListDepartments(context.Background(), 0)
	if err != nil || got == nil {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestJSAPISignature(t *testing.T) {
	got := JSAPISignature("sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90e5h", "Wm3WZYTPz0wzccnW", "http://mp.weixin.qq.com?params=value", 1414587457)
	if got == "" || len(got) != 40 {
		t.Fatalf("sig %s", got)
	}
}
