# wecom

Reusable WeCom (企业微信) self-built application client for Go. Stdlib only. Split so hosts import only what they need.

| Package | Import | What it is |
|---|---|---|
| core | `github.com/yibaiba/wecom` | `Client` (token + HTTP), `Identity`, `Production` / `Sandbox` login exchange |
| login pages | `github.com/yibaiba/wecom/login` | WWLogin panel, sandbox picker, phone QR |
| contacts | `github.com/yibaiba/wecom/contact` | members, departments, tags |
| application | `github.com/yibaiba/wecom/agent` | app settings, menu, workbench (template + per-user / batch data) |
| messages | `github.com/yibaiba/wecom/message` | send / recall / app group chat |
| media | `github.com/yibaiba/wecom/media` | upload / download |
| JS-SDK | `github.com/yibaiba/wecom/jsapi` | tickets + signature |
| callbacks | `github.com/yibaiba/wecom/callback` | receive-message encrypt / decrypt |

```bash
go get github.com/yibaiba/wecom@latest
```

```go
import (
    "github.com/yibaiba/wecom"
    "github.com/yibaiba/wecom/agent"
    "github.com/yibaiba/wecom/contact"
    "github.com/yibaiba/wecom/login"
    "github.com/yibaiba/wecom/message"
)

app := &wecom.Client{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}}

ident, err := wecom.Production{CorpID: corpID, AgentID: agentID, Secret: secret, HTTP: wecom.HTTPDoer{}}.Exchange(ctx, code)

ag, err := agent.New(app).Get(ctx)
users, err := contact.New(app).ListDepartmentUsers(ctx, 1)
_, err = message.New(app).SendText(ctx, "userid", "hello")
```

OA products (approval, calendar, meetings, WeDrive) are separate WeCom products and stay out of this module.

Empty Identity fields mean the corp app or the member did not grant them. `user_ticket` is short-lived; do not persist it.
