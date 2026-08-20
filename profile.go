package wecom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type userGetResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	UserID      string `json:"userid"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	BizMail     string `json:"biz_mail"`
	Mobile      string `json:"mobile"`
	Avatar      string `json:"avatar"`
	ThumbAvatar string `json:"thumb_avatar"`
	Gender      string `json:"gender"`
}

type userDetailResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	UserID  string `json:"userid"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Avatar  string `json:"avatar"`
	QRCode  string `json:"qr_code"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	BizMail string `json:"biz_mail"`
}

func (p Production) directoryProfile(ctx context.Context, token, userid string) (Identity, error) {
	u := getUserURL + "?access_token=" + url.QueryEscape(token) + "&userid=" + url.QueryEscape(userid)
	body, err := p.HTTP.Get(ctx, u)
	if err != nil {
		return Identity{}, fmt.Errorf("wecom user get: %w", err)
	}
	var out userGetResp
	if err := json.Unmarshal(body, &out); err != nil {
		return Identity{}, fmt.Errorf("wecom user get decode: %w", err)
	}
	if out.ErrCode != 0 {
		return Identity{}, fmt.Errorf("wecom user get error")
	}
	return Identity{
		UserID:   userid,
		Name:     strings.TrimSpace(out.Name),
		Alias:    strings.TrimSpace(out.Alias),
		Position: strings.TrimSpace(out.Position),
		Email:    firstNonEmpty(out.Email, out.BizMail),
		Mobile:   strings.TrimSpace(out.Mobile),
		Avatar:   firstNonEmpty(out.Avatar, out.ThumbAvatar),
		Gender:   strings.TrimSpace(out.Gender),
	}, nil
}

func (p Production) mergePrivateDetail(ctx context.Context, token, ticket string, ident Identity) Identity {
	poster, ok := p.HTTP.(poster)
	if !ok || strings.TrimSpace(ticket) == "" {
		return ident
	}
	payload, err := json.Marshal(map[string]string{"user_ticket": ticket})
	if err != nil {
		return ident
	}
	u := getUserDetailURL + "?access_token=" + url.QueryEscape(token)
	body, err := poster.Post(ctx, u, payload)
	if err != nil {
		return ident
	}
	var out userDetailResp
	if err := json.Unmarshal(body, &out); err != nil || out.ErrCode != 0 {
		return ident
	}
	ident.Name = firstNonEmpty(out.Name, ident.Name)
	ident.Email = firstNonEmpty(out.Email, out.BizMail, ident.Email)
	ident.Avatar = firstNonEmpty(out.Avatar, ident.Avatar)
	ident.Mobile = firstNonEmpty(out.Mobile, ident.Mobile)
	ident.Gender = firstNonEmpty(out.Gender, ident.Gender)
	return ident
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}
