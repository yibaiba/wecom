package message

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yibaiba/wecom"
)

func TestSendText(t *testing.T) {
	var sent map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/message/send":
			_ = json.NewDecoder(r.Body).Decode(&sent)
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "msgid": "m1"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	api := New(&wecom.Client{AgentID: 1000036, CorpID: "ww", Secret: "s", HTTP: wecom.HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL})
	res, err := api.SendText(context.Background(), "zhangsan", "hello")
	if err != nil || res.MsgID != "m1" || sent["msgtype"] != "text" {
		t.Fatalf("%v %+v %+v", err, res, sent)
	}
}
