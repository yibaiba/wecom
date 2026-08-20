package wecom

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type userResp struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	UserID         string `json:"UserId"`
	UserIDAlt      string `json:"userid"`
	DeviceID       string `json:"DeviceId"`
	OpenUserID     string `json:"open_userid"`
	UserTicket     string `json:"user_ticket"`
	UserDocTicket  string `json:"user_doc_ticket"`
	OpenID         string `json:"OpenId"`
	OpenIDAlt      string `json:"openid"`
	ExternalUserID string `json:"external_userid"`
}

func (u userResp) userid() string {
	if strings.TrimSpace(u.UserID) != "" {
		return strings.TrimSpace(u.UserID)
	}
	return strings.TrimSpace(u.UserIDAlt)
}

func (u userResp) openid() string {
	return firstNonEmpty(u.OpenID, u.OpenIDAlt)
}

// Exchange redeems a WeCom code for member identity. It maps every field
// user/get, auth/getuserinfo, and getuserdetail return. snsapi_privateinfo
// codes overlay sensitive fields via getuserdetail. Non-members return OpenID.
func (p Production) Exchange(ctx context.Context, code string) (Identity, error) {
	if p.HTTP == nil {
		return Identity{}, fmt.Errorf("wecom http client is not configured")
	}
	app := p.App()
	info, err := p.userFromCode(ctx, app, code)
	if err != nil {
		return Identity{}, err
	}
	return p.identityFromAuth(ctx, app, info)
}

func (p Production) identityFromAuth(ctx context.Context, app *Client, info userResp) (Identity, error) {
	userid := info.userid()
	if userid == "" {
		return applyAuthFields(Identity{}, info), nil
	}
	ident, err := p.directoryProfile(ctx, app, userid)
	if err != nil {
		if !isDirectoryScopeErr(err) {
			return Identity{}, err
		}
		ident = Identity{UserID: userid}
	}
	ident = applyAuthFields(ident, info)
	if info.UserTicket != "" {
		ident = p.mergePrivateDetail(ctx, app, info.UserTicket, ident)
	}
	if ident.Name == "" {
		ident.Name = ident.UserID
	}
	return ident, nil
}

func applyAuthFields(ident Identity, info userResp) Identity {
	ident.DeviceID = firstNonEmpty(info.DeviceID, ident.DeviceID)
	ident.OpenUserID = firstNonEmpty(ident.OpenUserID, info.OpenUserID)
	ident.UserTicket = firstNonEmpty(info.UserTicket, ident.UserTicket)
	ident.UserDocTicket = firstNonEmpty(info.UserDocTicket, ident.UserDocTicket)
	ident.OpenID = firstNonEmpty(info.openid(), ident.OpenID)
	ident.ExternalUserID = firstNonEmpty(info.ExternalUserID, ident.ExternalUserID)
	if ident.UserID == "" {
		ident.UserID = info.userid()
	}
	return ident
}

func (p Production) userFromCode(ctx context.Context, app *Client, code string) (userResp, error) {
	q := url.Values{}
	q.Set("code", code)
	var out userResp
	if err := app.GetJSON(ctx, "/cgi-bin/auth/getuserinfo", q, &out); err != nil {
		return userResp{}, err
	}
	if out.userid() == "" && out.openid() == "" && strings.TrimSpace(out.ExternalUserID) == "" {
		return userResp{}, Error{Msg: "empty userinfo"}
	}
	return out, nil
}
