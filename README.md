# wecom

Reusable WeCom (企业微信) login for Go. Stdlib only. Split so host apps import only what they need.

| Package | Import | What it is |
|---|---|---|
| adapter | `github.com/yibaiba/wecom` | `Production` / `Sandbox` code exchange. `Identity` maps every field `user/get`, `getuserinfo`, and `auth/getuserdetail` return. |
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
```

`Exchange` calls gettoken → getuserinfo → directory `user/get`. When the code has `user_ticket` (`snsapi_privateinfo`), it overlays `auth/getuserdetail`. Empty fields mean the corp app or the member did not grant them. New self-built apps usually need OAuth for avatar, gender, mobile, email, biz_mail, qr_code, and address.

Non-members (no `userid`) skip `user/get` and return `OpenID` / `ExternalUserID`.

Phone QR for unknown users: host your own status/continue routes and pass those paths into `login.WritePhoneQRPage`. Session cookies and employee tables stay in the host app.

## Identity

| Field | Source | Notes |
|---|---|---|
| UserID, Name, Alias, Position | `user/get` | Name falls back to UserID |
| ExternalPosition | `user/get` | Customer-facing title |
| Email | `user/get` then `getuserdetail` | Personal email; falls back to BizMail |
| BizMail | `user/get` then `getuserdetail` | Enterprise mailbox |
| Mobile, Telephone, Address | `user/get` / `getuserdetail` | Telephone is landline |
| Gender | `user/get` / `getuserdetail` | `0` undefined, `1` male, `2` female |
| Avatar, ThumbAvatar | `user/get` / `getuserdetail` | Avatar falls back to thumb |
| QRCode | `user/get` / `getuserdetail` | Personal contact QR URL |
| Status | `user/get` | `1` active, `2` disabled, `4` inactive, `5` quit |
| MainDepartment, Department, Order, IsLeaderInDept, DirectLeader | `user/get` | Org chart |
| ExtAttr | `user/get` | Custom text / web / miniprogram attrs |
| ExternalProfile | `user/get` | Corp name, video account, external attrs |
| OpenUserID | `user/get` / getuserinfo | Service-provider global id |
| DeviceID, UserTicket, UserDocTicket, OpenID, ExternalUserID | getuserinfo | Ticket is short-lived; do not persist |
