package wecom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

var (
	getTokenURL      = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	getUserInfoURL   = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
	getUserURL       = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	getUserDetailURL = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserdetail"
)

type tokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

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
// user/get, getuserinfo, and getuserdetail return. snsapi_privateinfo codes
// overlay sensitive fields via getuserdetail. Non-members return OpenID.
func (p Production) Exchange(ctx context.Context, code string) (Identity, error) {
	if p.HTTP == nil {
		return Identity{}, fmt.Errorf("wecom http client is not configured")
	}
	token, err := p.corpToken(ctx)
	if err != nil {
		return Identity{}, err
	}
	info, err := p.userFromCode(ctx, token, code)
	if err != nil {
		return Identity{}, err
	}
	return p.identityFromAuth(ctx, token, info)
}

func (p Production) identityFromAuth(ctx context.Context, token string, info userResp) (Identity, error) {
	userid := info.userid()
	if userid == "" {
		return applyAuthFields(Identity{}, info), nil
	}
	ident, err := p.directoryProfile(ctx, token, userid)
	if err != nil {
		return Identity{}, err
	}
	ident = applyAuthFields(ident, info)
	if info.UserTicket != "" {
		ident = p.mergePrivateDetail(ctx, token, info.UserTicket, ident)
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

func (p Production) corpToken(ctx context.Context) (string, error) {
	u := getTokenURL + "?corpid=" + url.QueryEscape(p.CorpID) + "&corpsecret=" + url.QueryEscape(p.Secret)
	body, err := p.HTTP.Get(ctx, u)
	if err != nil {
		return "", fmt.Errorf("wecom token: %w", err)
	}
	var out tokenResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("wecom token decode: %w", err)
	}
	if out.ErrCode != 0 || out.AccessToken == "" {
		return "", fmt.Errorf("wecom token error")
	}
	return out.AccessToken, nil
}

func (p Production) userFromCode(ctx context.Context, token, code string) (userResp, error) {
	u := getUserInfoURL + "?access_token=" + url.QueryEscape(token) + "&code=" + url.QueryEscape(code)
	body, err := p.HTTP.Get(ctx, u)
	if err != nil {
		return userResp{}, fmt.Errorf("wecom userinfo: %w", err)
	}
	var out userResp
	if err := json.Unmarshal(body, &out); err != nil {
		return userResp{}, fmt.Errorf("wecom userinfo decode: %w", err)
	}
	if out.ErrCode != 0 {
		return userResp{}, fmt.Errorf("wecom userinfo error")
	}
	if out.userid() == "" && out.openid() == "" && strings.TrimSpace(out.ExternalUserID) == "" {
		return userResp{}, fmt.Errorf("wecom userinfo error")
	}
	return out, nil
}
