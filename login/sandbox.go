package login

import (
	"html"
	"net/http"
	"net/url"
	"strings"

	"github.com/yibaiba/wecom"
)

// WriteSandboxPicker renders the explicit sandbox identity list.
// callbackPath is the host callback route (for example CallbackPath).
func WriteSandboxPicker(w http.ResponseWriter, state, callbackPath string, users []wecom.Identity) {
	if callbackPath == "" {
		callbackPath = CallbackPath
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", pagePolicy())
	var b strings.Builder
	b.WriteString("<html><body><h1>Sandbox WeCom</h1><ul>")
	for _, u := range users {
		href := callbackPath + "?state=" + url.QueryEscape(state) + "&code=" + url.QueryEscape(u.UserID)
		b.WriteString("<li><a href=\"" + html.EscapeString(href) + "\">" + html.EscapeString(u.Name) + "</a></li>")
	}
	b.WriteString("</ul></body></html>")
	_, _ = w.Write([]byte(b.String()))
}
