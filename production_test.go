package wecom

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testProduction(t *testing.T, h http.HandlerFunc) Production {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return Production{
		CorpID:  "ww",
		Secret:  "s",
		HTTP:    HTTPDoer{Client: srv.Client()},
		BaseURL: srv.URL,
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	_ = json.NewEncoder(w).Encode(v)
}

func TestExchangeReadsEmailAndAvatar(t *testing.T) {
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeJSON(w, tokenResp{AccessToken: "tok", ErrCode: 0})
		case "/cgi-bin/user/getuserinfo":
			t.Errorf("legacy getuserinfo path")
			http.NotFound(w, r)
		case "/cgi-bin/auth/getuserinfo":
			writeJSON(w, map[string]any{"UserId": "zhangsan", "errcode": 0, "user_ticket": "ticket-1"})
		case "/cgi-bin/user/get":
			writeJSON(w, userGetResp{
				ErrCode: 0, UserID: "zhangsan", Name: "张三", Email: "dir@example.com", Avatar: "https://img.example/dir.png",
			})
		case "/cgi-bin/auth/getuserdetail":
			writeJSON(w, userDetailResp{
				ErrCode: 0, Name: "张三", Email: "priv@example.com", Avatar: "https://img.example/priv.png", Mobile: "13800000000",
			})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := p.Exchange(context.Background(), "code-1")
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
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeJSON(w, tokenResp{AccessToken: "tok"})
		case "/cgi-bin/auth/getuserinfo":
			writeJSON(w, map[string]any{"userid": "lisi", "errcode": 0})
		case "/cgi-bin/user/get":
			writeJSON(w, userGetResp{Name: "李四", BizMail: "lisi@corp.com", ThumbAvatar: "https://img.example/lisi.png"})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := p.Exchange(context.Background(), "code-2")
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != "lisi" || got.Email != "" || got.Avatar != "https://img.example/lisi.png" {
		t.Fatalf("identity %+v", got)
	}
	if got.BizMail != "lisi@corp.com" || got.ThumbAvatar != "https://img.example/lisi.png" {
		t.Fatalf("biz/thumb %+v", got)
	}
}

func TestExchangeMapsOfficialMemberPayload(t *testing.T) {
	userGetCalled := 0
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeJSON(w, tokenResp{AccessToken: "tok"})
		case "/cgi-bin/auth/getuserinfo":
			writeJSON(w, map[string]any{
				"errcode": 0, "UserId": "zhangsan", "DeviceId": "dev-1",
				"open_userid": "open-from-info", "user_ticket": "ticket-1",
				"user_doc_ticket": "doc-1",
			})
		case "/cgi-bin/user/get":
			userGetCalled++
			_, _ = w.Write([]byte(officialUserGetJSON))
		case "/cgi-bin/auth/getuserdetail":
			writeJSON(w, map[string]any{
				"errcode": 0, "userid": "zhangsan", "gender": "1",
				"avatar":  "https://img.example/priv.png",
				"qr_code": "https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=priv",
				"mobile":  "13900000000", "email": "priv@example.com",
				"biz_mail": "priv@qy.wecom.work", "address": "深圳市南山区",
			})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := p.Exchange(context.Background(), "code-full")
	if err != nil {
		t.Fatal(err)
	}
	if userGetCalled != 1 {
		t.Fatalf("user/get calls %d", userGetCalled)
	}
	assertFullMember(t, got)
}

func TestExchangeVisitorSkipsDirectory(t *testing.T) {
	userGetCalled := 0
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeJSON(w, tokenResp{AccessToken: "tok"})
		case "/cgi-bin/auth/getuserinfo":
			writeJSON(w, map[string]any{
				"errcode": 0, "openid": "oid-1", "external_userid": "ext-1", "DeviceId": "dev-ext",
			})
		case "/cgi-bin/user/get":
			userGetCalled++
			http.Error(w, "should not call user/get", http.StatusInternalServerError)
		default:
			http.NotFound(w, r)
		}
	})
	got, err := p.Exchange(context.Background(), "code-ext")
	if err != nil {
		t.Fatal(err)
	}
	if userGetCalled != 0 {
		t.Fatalf("user/get calls %d", userGetCalled)
	}
	if got.UserID != "" || got.OpenID != "oid-1" || got.ExternalUserID != "ext-1" || got.DeviceID != "dev-ext" {
		t.Fatalf("visitor %+v", got)
	}
}

func TestExchangeContinuesWhenUserOutOfScope(t *testing.T) {
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			writeJSON(w, tokenResp{AccessToken: "tok"})
		case "/cgi-bin/auth/getuserinfo":
			writeJSON(w, map[string]any{"errcode": 0, "userid": "hidden", "user_ticket": "ticket-1"})
		case "/cgi-bin/user/get":
			writeJSON(w, map[string]any{"errcode": 81013, "errmsg": "user out of scope"})
		case "/cgi-bin/auth/getuserdetail":
			writeJSON(w, map[string]any{"errcode": 0, "userid": "hidden", "name": "隐身", "mobile": "13800000000"})
		default:
			http.NotFound(w, r)
		}
	})
	got, err := p.Exchange(context.Background(), "code-hidden")
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != "hidden" || got.Name != "隐身" || got.Mobile != "13800000000" {
		t.Fatalf("identity %+v", got)
	}
}

func TestExchangeReturnsAPIError(t *testing.T) {
	p := testProduction(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			writeJSON(w, map[string]any{"errcode": 40013, "errmsg": "invalid corpid"})
			return
		}
		http.NotFound(w, r)
	})
	_, err := p.Exchange(context.Background(), "code")
	var e Error
	if !errors.As(err, &e) || e.Code != 40013 {
		t.Fatalf("err %v", err)
	}
}

func assertFullMember(t *testing.T, got Identity) {
	t.Helper()
	if got.UserID != "zhangsan" || got.Name != "张三" || got.Alias != "jackzhang" {
		t.Fatalf("ids %+v", got)
	}
	if got.Position != "后台工程师" || got.ExternalPosition != "产品经理" {
		t.Fatalf("position %+v", got)
	}
	if got.Email != "priv@example.com" || got.BizMail != "priv@qy.wecom.work" {
		t.Fatalf("mail %+v", got)
	}
	if got.Mobile != "13900000000" || got.Telephone != "020-123456" || got.Address != "深圳市南山区" {
		t.Fatalf("contact %+v", got)
	}
	if got.Gender != "1" || got.Avatar != "https://img.example/priv.png" || got.ThumbAvatar == "" {
		t.Fatalf("portrait %+v", got)
	}
	if got.QRCode != "https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=priv" {
		t.Fatalf("qr %s", got.QRCode)
	}
	if got.Status != 1 || got.MainDepartment != 1 {
		t.Fatalf("status %+v", got)
	}
	if len(got.Department) != 2 || got.Department[0] != 1 || got.Order[1] != 2 || got.IsLeaderInDept[0] != 1 {
		t.Fatalf("dept %+v", got)
	}
	if len(got.DirectLeader) != 1 || got.DirectLeader[0] != "lisi" {
		t.Fatalf("leader %+v", got)
	}
	if got.OpenUserID != "wo-open-userid" || got.DeviceID != "dev-1" || got.UserTicket != "ticket-1" || got.UserDocTicket != "doc-1" {
		t.Fatalf("auth %+v", got)
	}
	if len(got.ExtAttr) != 2 || got.ExtAttr[0].Text.Value != "文本" || got.ExtAttr[1].Web.URL != "http://www.test.com" {
		t.Fatalf("extattr %+v", got.ExtAttr)
	}
	if got.ExternalProfile.ExternalCorpName != "企业简称" || got.ExternalProfile.WechatChannels.Nickname != "视频号名称" {
		t.Fatalf("external %+v", got.ExternalProfile)
	}
	if len(got.ExternalProfile.ExternalAttr) != 3 || got.ExternalProfile.ExternalAttr[2].Miniprogram.AppID == "" {
		t.Fatalf("external attr %+v", got.ExternalProfile.ExternalAttr)
	}
}

const officialUserGetJSON = `{
	"errcode": 0,
	"errmsg": "ok",
	"userid": "zhangsan",
	"name": "张三",
	"department": [1, 2],
	"order": [1, 2],
	"position": "后台工程师",
	"mobile": "13800000000",
	"gender": "1",
	"email": "zhangsan@qq.com",
	"biz_mail": "zhangsan@tencent.com",
	"is_leader_in_dept": [1, 0],
	"direct_leader": ["lisi"],
	"avatar": "http://wx.qlogo.cn/mmopen/dir/0",
	"thumb_avatar": "http://wx.qlogo.cn/mmopen/dir/100",
	"telephone": "020-123456",
	"alias": "jackzhang",
	"address": "广州市海珠区新港中路",
	"open_userid": "wo-open-userid",
	"main_department": 1,
	"extattr": {
		"attrs": [
			{"type": 0, "name": "文本名称", "text": {"value": "文本"}},
			{"type": 1, "name": "网页名称", "web": {"url": "http://www.test.com", "title": "标题"}}
		]
	},
	"status": 1,
	"qr_code": "https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=dir",
	"external_position": "产品经理",
	"external_profile": {
		"external_corp_name": "企业简称",
		"wechat_channels": {"nickname": "视频号名称", "status": 1},
		"external_attr": [
			{"type": 0, "name": "文本名称", "text": {"value": "文本"}},
			{"type": 1, "name": "网页名称", "web": {"url": "http://www.test.com", "title": "标题"}},
			{"type": 2, "name": "测试app", "miniprogram": {"appid": "wx8bd80126147dFAKE", "pagepath": "/index", "title": "my miniprogram"}}
		]
	}
}`
