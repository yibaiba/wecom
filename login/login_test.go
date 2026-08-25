package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteLoginPanel(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteLoginPanel(rec, App{CorpID: "wwpanel", AgentID: 1000002}, "abc", "https://app.example.test/auth/wecom/callback")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "createWWLoginPanel") || !strings.Contains(body, "wwpanel") {
		t.Fatalf("panel %s", body)
	}
	if strings.Contains(body, "qrConnect") || strings.Contains(body, "lockFrame") {
		t.Fatal("panel must stay on the WWLogin component")
	}
}

func TestWritePhoneQRPageUsesHostPaths(t *testing.T) {
	rec := httptest.NewRecorder()
	WritePhoneQRPage(rec, "https://open.weixin.qq.com/connect/oauth2/authorize?scope=snsapi_privateinfo", "/auth/wecom/status", "/auth/wecom/continue")
	body := rec.Body.String()
	if !strings.Contains(body, "请使用手机企业微信扫码") || !strings.Contains(body, "snsapi_privateinfo") {
		t.Fatalf("qr page %s", body)
	}
	if !strings.Contains(body, "/auth/wecom/status") || !strings.Contains(body, "/auth/wecom/continue") {
		t.Fatal("qr page must use host poll paths")
	}
	if strings.Contains(body, "cdn.jsdelivr.net") || strings.Contains(body, "QRCode.toCanvas") {
		t.Fatal("qr page must not load QR from external CDN")
	}
	if !strings.Contains(body, "data:image/png;base64,") {
		t.Fatal("qr page must embed a server-rendered PNG")
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("cache-control %q", rec.Header().Get("Cache-Control"))
	}
}

func TestIsWxWorkUA(t *testing.T) {
	if !IsWxWorkUA("Mozilla/5.0 wxwork") || IsWxWorkUA("Mozilla/5.0 Chrome") {
		t.Fatal("wxwork detection")
	}
}
