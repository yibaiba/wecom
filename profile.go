package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type userGetResp struct {
	ErrCode          int            `json:"errcode"`
	ErrMsg           string         `json:"errmsg"`
	UserID           string         `json:"userid"`
	Name             string         `json:"name"`
	Alias            string         `json:"alias"`
	Position         string         `json:"position"`
	ExternalPosition string         `json:"external_position"`
	Email            string         `json:"email"`
	BizMail          string         `json:"biz_mail"`
	Mobile           string         `json:"mobile"`
	Telephone        string         `json:"telephone"`
	Address          string         `json:"address"`
	Avatar           string         `json:"avatar"`
	ThumbAvatar      string         `json:"thumb_avatar"`
	QRCode           string         `json:"qr_code"`
	Gender           wecomString    `json:"gender"`
	Status           int            `json:"status"`
	MainDepartment   int            `json:"main_department"`
	Department       []int          `json:"department"`
	Order            []int          `json:"order"`
	IsLeaderInDept   []int          `json:"is_leader_in_dept"`
	DirectLeader     []string       `json:"direct_leader"`
	OpenUserID       string         `json:"open_userid"`
	ExtAttr          extAttrWrap    `json:"extattr"`
	ExternalProfile  extProfileWire `json:"external_profile"`
}

type userDetailResp struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	UserID  string      `json:"userid"`
	Name    string      `json:"name"`
	Gender  wecomString `json:"gender"`
	Avatar  string      `json:"avatar"`
	QRCode  string      `json:"qr_code"`
	Mobile  string      `json:"mobile"`
	Email   string      `json:"email"`
	BizMail string      `json:"biz_mail"`
	Address string      `json:"address"`
}

type extAttrWrap struct {
	Attrs []extAttrWire `json:"attrs"`
}

type extAttrWire struct {
	Type        int            `json:"type"`
	Name        string         `json:"name"`
	Text        ExtAttrText    `json:"text"`
	Web         extAttrWebWire `json:"web"`
	Miniprogram extAttrAppWire `json:"miniprogram"`
}

type extAttrWebWire struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type extAttrAppWire struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
	Title    string `json:"title"`
}

type extProfileWire struct {
	ExternalCorpName string            `json:"external_corp_name"`
	WechatChannels   wechatChannelWire `json:"wechat_channels"`
	ExternalAttr     []extAttrWire     `json:"external_attr"`
}

type wechatChannelWire struct {
	Nickname string `json:"nickname"`
	Status   int    `json:"status"`
}

// wecomString accepts JSON string or number (gender is documented as string).
type wecomString string

func (s *wecomString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*s = ""
		return nil
	}
	if b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		*s = wecomString(v)
		return nil
	}
	*s = wecomString(string(b))
	return nil
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
	return identityFromUserGet(userid, out), nil
}

func identityFromUserGet(userid string, out userGetResp) Identity {
	email := strings.TrimSpace(out.Email)
	biz := strings.TrimSpace(out.BizMail)
	avatar := strings.TrimSpace(out.Avatar)
	thumb := strings.TrimSpace(out.ThumbAvatar)
	return Identity{
		UserID:           userid,
		Name:             strings.TrimSpace(out.Name),
		Alias:            strings.TrimSpace(out.Alias),
		Position:         strings.TrimSpace(out.Position),
		ExternalPosition: strings.TrimSpace(out.ExternalPosition),
		Email:            firstNonEmpty(email, biz),
		BizMail:          biz,
		Mobile:           strings.TrimSpace(out.Mobile),
		Telephone:        strings.TrimSpace(out.Telephone),
		Address:          strings.TrimSpace(out.Address),
		Gender:           strings.TrimSpace(string(out.Gender)),
		Avatar:           firstNonEmpty(avatar, thumb),
		ThumbAvatar:      thumb,
		QRCode:           strings.TrimSpace(out.QRCode),
		Status:           out.Status,
		MainDepartment:   out.MainDepartment,
		Department:       out.Department,
		Order:            out.Order,
		IsLeaderInDept:   out.IsLeaderInDept,
		DirectLeader:     out.DirectLeader,
		OpenUserID:       strings.TrimSpace(out.OpenUserID),
		ExtAttr:          extAttrsFromWire(out.ExtAttr.Attrs),
		ExternalProfile:  externalProfileFromWire(out.ExternalProfile),
	}
}

// ParseDirectoryUser maps a user/get JSON body onto Identity.
func ParseDirectoryUser(userid string, body []byte) (Identity, error) {
	var out userGetResp
	if err := json.Unmarshal(body, &out); err != nil {
		return Identity{}, fmt.Errorf("wecom user get decode: %w", err)
	}
	if userid == "" {
		userid = out.UserID
	}
	return identityFromUserGet(userid, out), nil
}

// ParseDirectoryUsers maps a user/list JSON body onto identities.
func ParseDirectoryUsers(body []byte) ([]Identity, error) {
	var out struct {
		UserList []userGetResp `json:"userlist"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("wecom user list decode: %w", err)
	}
	users := make([]Identity, 0, len(out.UserList))
	for _, row := range out.UserList {
		users = append(users, identityFromUserGet(row.UserID, row))
	}
	return users, nil
}

func (p Production) mergePrivateDetail(ctx context.Context, token, ticket string, ident Identity) Identity {
	if p.HTTP == nil || strings.TrimSpace(ticket) == "" {
		return ident
	}
	payload, err := json.Marshal(map[string]string{"user_ticket": ticket})
	if err != nil {
		return ident
	}
	u := getUserDetailURL + "?access_token=" + url.QueryEscape(token)
	body, err := p.HTTP.Post(ctx, u, payload)
	if err != nil {
		return ident
	}
	var out userDetailResp
	if err := json.Unmarshal(body, &out); err != nil || out.ErrCode != 0 {
		return ident
	}
	return overlayPrivateDetail(ident, out)
}

func overlayPrivateDetail(ident Identity, out userDetailResp) Identity {
	ident.Name = firstNonEmpty(out.Name, ident.Name)
	ident.Gender = firstNonEmpty(string(out.Gender), ident.Gender)
	ident.Avatar = firstNonEmpty(out.Avatar, ident.Avatar)
	ident.QRCode = firstNonEmpty(out.QRCode, ident.QRCode)
	ident.Mobile = firstNonEmpty(out.Mobile, ident.Mobile)
	ident.Email = firstNonEmpty(out.Email, out.BizMail, ident.Email)
	ident.BizMail = firstNonEmpty(out.BizMail, ident.BizMail)
	ident.Address = firstNonEmpty(out.Address, ident.Address)
	if uid := strings.TrimSpace(out.UserID); uid != "" {
		ident.UserID = firstNonEmpty(ident.UserID, uid)
	}
	return ident
}

func extAttrsFromWire(attrs []extAttrWire) []ExtAttr {
	if len(attrs) == 0 {
		return nil
	}
	out := make([]ExtAttr, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, ExtAttr{
			Type: a.Type,
			Name: strings.TrimSpace(a.Name),
			Text: ExtAttrText{Value: strings.TrimSpace(a.Text.Value)},
			Web: ExtAttrWeb{
				URL:   strings.TrimSpace(a.Web.URL),
				Title: strings.TrimSpace(a.Web.Title),
			},
			Miniprogram: ExtAttrMiniprogram{
				AppID:    strings.TrimSpace(a.Miniprogram.AppID),
				PagePath: strings.TrimSpace(a.Miniprogram.PagePath),
				Title:    strings.TrimSpace(a.Miniprogram.Title),
			},
		})
	}
	return out
}

func externalProfileFromWire(p extProfileWire) ExternalProfile {
	return ExternalProfile{
		ExternalCorpName: strings.TrimSpace(p.ExternalCorpName),
		WechatChannels: WechatChannels{
			Nickname: strings.TrimSpace(p.WechatChannels.Nickname),
			Status:   p.WechatChannels.Status,
		},
		ExternalAttr: extAttrsFromWire(p.ExternalAttr),
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}
