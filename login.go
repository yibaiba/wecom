package wecom

import (
	"strings"
)

// Default login paths. Host apps may use other routes; pass those paths into
// WriteSandboxPicker and WritePhoneQRPage.
const (
	StartPath          = "/login/wecom"
	CallbackPath       = "/login/wecom/callback"
	EnrollCallbackPath = "/login/wecom/enroll/callback"
	EnrollStatusPath   = "/login/wecom/enroll/status"
	EnrollContinuePath = "/login/wecom/enroll/continue"
)

// LoginApp is the WeCom corp application used for browser login.
type LoginApp struct {
	CorpID  string
	AgentID int
}

// Configured reports whether corp id and agent id are present.
func (a LoginApp) Configured() bool {
	return a.CorpID != "" && a.AgentID != 0
}

// PublicURL joins a public origin with a login path.
func PublicURL(publicBase, path string) string {
	return strings.TrimRight(publicBase, "/") + path
}

// IsWxWorkUA reports the official WeCom built-in browser.
func IsWxWorkUA(ua string) bool {
	return strings.Contains(strings.ToLower(ua), "wxwork")
}

// WebviewURL is silent OAuth2 for the WeCom built-in browser.
func (a LoginApp) WebviewURL(callback, state string) string {
	return WxWorkAuthURL(a.CorpID, a.AgentID, callback, state)
}

// PrivateInfoURL is the phone-scan OAuth2 URL for unknown members.
func (a LoginApp) PrivateInfoURL(callback, state string) string {
	return WxWorkPrivateInfoURL(a.CorpID, a.AgentID, callback, state)
}
