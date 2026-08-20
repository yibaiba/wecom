# wecom

Reusable WeCom (企业微信) **self-built application** client for Go. Stdlib only.

A corp app can log members in, read the address book in its visible range, send messages, manage media/menu, and receive callbacks. This module covers those application APIs. OA products (approval, calendar, meetings, WeDrive, …) stay out of this package — they are separate WeCom products, not what every 应用 can do.

| Package | Import | What it is |
|---|---|---|
| app + login exchange | `github.com/yibaiba/wecom` | `Client` for application APIs. `Production` / `Sandbox` for OAuth `Exchange`. |
| login pages | `github.com/yibaiba/wecom/login` | WWLogin panel, sandbox picker, phone QR |

```bash
go get github.com/yibaiba/wecom@latest
```

```go
cli := &wecom.Client{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}}

ident, err := wecom.Production{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}}.Exchange(ctx, code)

ag, err := cli.GetAgent(ctx)                 // visible users / depts / tags
users, err := cli.ListDepartmentUsers(ctx, 1)
_, err = cli.SendText(ctx, "userid", "hello")
media, err := cli.UploadMediaBytes(ctx, "file", "a.txt", data)
ticket, err := cli.JSAPITicket(ctx)
sig := wecom.JSAPISignature(ticket, nonce, url, ts)
```

Receive-message URL:

```go
cb := wecom.Callback{Token: token, EncodingAESKey: aesKey, CorpID: corpID}
echo, err := cb.VerifyURL(msgSig, timestamp, nonce, echostr)
plain, err := cb.Decrypt(msgSig, timestamp, nonce, body)
```

## Client surface

**Auth / identity** — `Production.Exchange`, `WxWorkAuthURL`, `WxWorkPrivateInfoURL` (see Identity below).

**Application** — `GetAgent`, `ListAgents`, `SetAgent`, `CreateMenu`, `GetMenu`, `DeleteMenu`, workbench template/data.

**Contacts** — `GetUser`, `ListDepartmentUsers`, `ListDepartmentUserDetails`, `CreateUser`, `UpdateUser`, `DeleteUser`, `BatchDeleteUsers`, `ListUserIDs`, `UserIDByMobile`, `UserIDByEmail`, `UserIDToOpenID`, `OpenIDToUserID`, `InviteUsers`, `AuthSuccess`, `JoinQRCode`.

**Departments / tags** — list/get/create/update/delete, tag members.

**Messages** — `Send` (text, image, voice, video, file, textcard, news, mpnews, markdown, miniprogram_notice, template_card JSON), `SendText`, `SendMarkdown`, `RecallMessage`, `UpdateTemplateCard`, app group chat create/update/get/send.

**Media** — `UploadMedia`, `UploadImage`, `GetMedia`.

**JS-SDK** — `JSAPITicket`, `AgentTicket`, `JSAPISignature`.

**Network** — `APIDomainIPs`, `CallbackIPs`.

**Callbacks** — `Callback.VerifyURL`, `Decrypt`, `Encrypt`.

Empty fields mean the corp app or the member did not grant them. New self-built apps usually need OAuth for avatar, gender, mobile, email, biz_mail, qr_code, and address.

## Identity

`Exchange` calls gettoken → getuserinfo → directory `user/get`. When the code has `user_ticket` (`snsapi_privateinfo`), it overlays `auth/getuserdetail`. Non-members skip `user/get` and return `OpenID` / `ExternalUserID`.

| Field | Source |
|---|---|
| UserID, Name, Alias, Position, ExternalPosition | `user/get` |
| Email, BizMail, Mobile, Telephone, Address, Gender, Avatar, ThumbAvatar, QRCode | `user/get` / `getuserdetail` |
| Status, org fields, ExtAttr, ExternalProfile | `user/get` |
| DeviceID, UserTicket, UserDocTicket, OpenID, ExternalUserID | getuserinfo |

Phone QR for unknown users: host status/continue routes and pass those paths into `login.WritePhoneQRPage`. Session cookies stay in the host app.
