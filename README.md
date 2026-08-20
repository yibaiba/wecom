# wecom

Reusable WeCom (企业微信) login module. Stdlib only. Host apps do not import enterprise-sso.

```bash
go get github.com/yibaiba/wecom@latest
```

```go
import "github.com/yibaiba/wecom"

app := wecom.LoginApp{CorpID: corpID, AgentID: agentID}
ex := wecom.Production{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}, RedirectURI: callback}

if wecom.IsWxWorkUA(r.UserAgent()) {
    http.Redirect(w, r, app.WebviewURL(callback, state), http.StatusFound)
    return
}
wecom.WriteLoginPanel(w, app, state, callback)
ident, err := ex.Exchange(r.Context(), code)
```

Phone QR for unknown users: host your own status/continue routes and pass those paths into `WritePhoneQRPage`. Session cookies and employee tables stay in the host app.
