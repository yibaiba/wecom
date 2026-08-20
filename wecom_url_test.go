package wecom

import (
	"strings"
	"testing"
)

func TestWxWorkAuthURL(t *testing.T) {
	got := WxWorkAuthURL("wwcorp", 12, "https://app.example.test/auth/wecom/callback", "st")
	if !strings.Contains(got, "https://open.weixin.qq.com/connect/oauth2/authorize") {
		t.Fatalf("url %s", got)
	}
	if !strings.HasSuffix(got, "#wechat_redirect") {
		t.Fatalf("missing fragment %s", got)
	}
	if !strings.Contains(got, "appid=wwcorp") || !strings.Contains(got, "agentid=12") {
		t.Fatalf("url %s", got)
	}
	if !strings.Contains(got, "scope=snsapi_base") {
		t.Fatalf("url %s", got)
	}
}

func TestWxWorkPrivateInfoURL(t *testing.T) {
	got := WxWorkPrivateInfoURL("wwcorp", 12, "https://app.example.test/auth/wecom/enroll/callback", "st")
	if !strings.Contains(got, "scope=snsapi_privateinfo") {
		t.Fatalf("url %s", got)
	}
	if !strings.Contains(got, "auth%2Fwecom%2Fenroll%2Fcallback") {
		t.Fatalf("url %s", got)
	}
}
