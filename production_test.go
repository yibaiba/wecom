package wecom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeReadsEmailAndAvatar(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(tokenResp{AccessToken: "tok", ErrCode: 0})
		case "/cgi-bin/user/getuserinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{"UserId": "zhangsan", "errcode": 0, "user_ticket": "ticket-1"})
		case "/cgi-bin/user/get":
			_ = json.NewEncoder(w).Encode(userGetResp{
				ErrCode: 0, UserID: "zhangsan", Name: "张三", Email: "dir@example.com", Avatar: "https://img.example/dir.png",
			})
		case "/cgi-bin/auth/getuserdetail":
			_ = json.NewEncoder(w).Encode(userDetailResp{
				ErrCode: 0, Name: "张三", Email: "priv@example.com", Avatar: "https://img.example/priv.png", Mobile: "13800000000",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	origToken, origInfo, origUser, origDetail := getTokenURL, getUserInfoURL, getUserURL, getUserDetailURL
	getTokenURL, getUserInfoURL, getUserURL, getUserDetailURL = srv.URL+"/cgi-bin/gettoken", srv.URL+"/cgi-bin/user/getuserinfo", srv.URL+"/cgi-bin/user/get", srv.URL+"/cgi-bin/auth/getuserdetail"
	defer func() {
		getTokenURL, getUserInfoURL, getUserURL, getUserDetailURL = origToken, origInfo, origUser, origDetail
	}()

	got, err := (Production{CorpID: "ww", Secret: "s", HTTP: HTTPDoer{Client: srv.Client()}}).Exchange(context.Background(), "code-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != "zhangsan" || got.Name != "张三" {
		t.Fatalf("identity %+v", got)
	}
	if got.Email != "priv@example.com" || got.Avatar != "https://img.example/priv.png" || got.Mobile != "13800000000" {
		t.Fatalf("profile %+v", got)
	}
}

func TestExchangeKeepsDirectoryProfileWithoutTicket(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_ = json.NewEncoder(w).Encode(tokenResp{AccessToken: "tok"})
		case "/cgi-bin/user/getuserinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{"userid": "lisi", "errcode": 0})
		case "/cgi-bin/user/get":
			_ = json.NewEncoder(w).Encode(userGetResp{Name: "李四", BizMail: "lisi@corp.com", ThumbAvatar: "https://img.example/lisi.png"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	origToken, origInfo, origUser := getTokenURL, getUserInfoURL, getUserURL
	getTokenURL, getUserInfoURL, getUserURL = srv.URL+"/cgi-bin/gettoken", srv.URL+"/cgi-bin/user/getuserinfo", srv.URL+"/cgi-bin/user/get"
	defer func() { getTokenURL, getUserInfoURL, getUserURL = origToken, origInfo, origUser }()

	got, err := (Production{CorpID: "ww", Secret: "s", HTTP: HTTPDoer{Client: srv.Client()}}).Exchange(context.Background(), "code-2")
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != "lisi" || got.Email != "lisi@corp.com" || got.Avatar != "https://img.example/lisi.png" {
		t.Fatalf("identity %+v", got)
	}
}
