package contact

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yibaiba/wecom"
)

func testAPI(t *testing.T, h http.HandlerFunc) *API {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return New(&wecom.Client{CorpID: "ww", Secret: "s", HTTP: wecom.HTTPDoer{Client: srv.Client()}, BaseURL: srv.URL})
}

func TestListDepartments(t *testing.T) {
	api := testAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/department/list":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errcode":    0,
				"department": []map[string]any{{"id": 1, "name": "Root", "parentid": 0}},
			})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := api.ListDepartments(context.Background(), 0)
	if err != nil || len(got) != 1 || got[0].Name != "Root" {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestGetUser(t *testing.T) {
	api := testAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok", "expires_in": 7200})
		case "/cgi-bin/user/get":
			_ = json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "userid": "zhangsan", "name": "张三", "email": "a@b.com"})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := api.GetUser(context.Background(), "zhangsan")
	if err != nil || got.Name != "张三" || got.Email != "a@b.com" {
		t.Fatalf("%v %+v", err, got)
	}
}
