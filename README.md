# wecom

Reusable WeCom (企业微信) login for Go. Stdlib only. Split so host apps import only what they need.

| Package | Import | What it is |
|---|---|---|
| adapter | `github.com/yibaiba/wecom` | `Production` / `Sandbox` code exchange. `Identity` includes userid, name, email, avatar, mobile when the WeCom app can read the address book; `snsapi_privateinfo` fills gaps via `getuserdetail`. |
| login pages | `github.com/yibaiba/wecom/login` | WWLogin panel, sandbox picker, phone QR, default routes |

```bash
go get github.com/yibaiba/wecom@latest
```

```go
import (
    "github.com/yibaiba/wecom"
    "github.com/yibaiba/wecom/login"
)

app := login.App{CorpID: corpID, AgentID: agentID}
ex := wecom.Production{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}, RedirectURI: callback}

if login.IsWxWorkUA(r.UserAgent()) {
    http.Redirect(w, r, app.WebviewURL(callback, state), http.StatusFound)
    return
}
login.WriteLoginPanel(w, app, state, callback)
ident, err := ex.Exchange(r.Context(), code)
// ident.Email, ident.Avatar, ident.Mobile, ident.Name, ident.UserID
```

Phone QR for unknown users: host your own status/continue routes and pass those paths into `login.WritePhoneQRPage`. Session cookies and employee tables stay in the host app.
