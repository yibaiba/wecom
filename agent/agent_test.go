package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yibaiba/wecom"
)

func TestGetAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/agent/get":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errcode": 0, "agentid": 1000036, "name": "SSO",
				"allow_userinfos": map[string]any{"user": []any{map[string]string{"userid": "zhangsan"}}},
				"allow_partys":    map[string]any{"partyid": []int{1}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	api := New(&wecom.Client{AgentID: 1000036, CorpID: "ww", Secret: "s", HTTP: wecom.HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL})
	ag, err := api.Get(context.Background())
	if err != nil || ag.Name != "SSO" || ag.AllowUserIDs[0] != "zhangsan" {
		t.Fatalf("%v %+v", err, ag)
	}
}
