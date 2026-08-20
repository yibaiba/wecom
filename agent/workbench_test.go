package agent

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yibaiba/wecom"
)

func workbenchAPI(t *testing.T, h http.HandlerFunc) *API {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return New(&wecom.Client{AgentID: 1000005, CorpID: "ww", Secret: "s", HTTP: wecom.HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL})
}

func TestBatchSetWorkbenchData(t *testing.T) {
	var got map[string]any
	api := workbenchAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/agent/batch_set_workbench_data":
			body, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(body, &got)
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "errmsg": "ok"})
		default:
			http.NotFound(w, r)
		}
	})
	err := api.BatchSetWorkbenchData(context.Background(), []string{"userid1", "userid2"}, WorkbenchData{
		Type:    "keydata",
		KeyData: &WorkbenchKeyData{Items: []WorkbenchKeyItem{{Key: "待审批", Data: "0"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got["agentid"].(float64) != 1000005 {
		t.Fatalf("agentid %+v", got)
	}
	list, _ := got["userid_list"].([]any)
	if len(list) != 2 {
		t.Fatalf("users %+v", got["userid_list"])
	}
}

func TestGetWorkbenchData(t *testing.T) {
	var req map[string]any
	api := workbenchAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/agent/get_workbench_data":
			_ = json.NewDecoder(r.Body).Decode(&req)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errcode": 0,
				"data": map[string]any{
					"type": "keydata",
					"keydata": map[string]any{
						"items": []any{map[string]any{"key": "待审批", "data": "2"}},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := api.GetWorkbenchData(context.Background(), "userid1")
	if err != nil {
		t.Fatal(err)
	}
	if req["userid"] != "userid1" || req["agentid"].(float64) != 1000005 {
		t.Fatalf("req %+v", req)
	}
	if got.Type != "keydata" || got.KeyData == nil || len(got.KeyData.Items) != 1 || got.KeyData.Items[0].Data != "2" {
		t.Fatalf("data %+v", got)
	}
}

func TestGetWorkbenchTemplatePostsJSON(t *testing.T) {
	var method string
	api := workbenchAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/agent/get_workbench_template":
			method = r.Method
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "type": "image"})
		default:
			http.NotFound(w, r)
		}
	})
	out, err := api.GetWorkbenchTemplate(context.Background())
	if err != nil || method != http.MethodPost || out["type"] != "image" {
		t.Fatalf("method %s out %+v err %v", method, out, err)
	}
}
