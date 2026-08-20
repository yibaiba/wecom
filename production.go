package wecom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type tokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

type userResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	UserID     string `json:"UserId"`
	DeviceID   string `json:"DeviceId"`
	OpenUserID string `json:"open_userid"`
}

type userGetResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Name    string `json:"name"`
	UserID  string `json:"userid"`
}

// HTTPDoer performs GET requests.
type HTTPDoer struct {
	Client *http.Client
}

// Get fetches a URL.
func (h HTTPDoer) Get(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	client := h.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}

// Exchange redeem the WeCom code for a userid and display name.
func (p Production) Exchange(ctx context.Context, code string) (Identity, error) {
	if p.HTTP == nil {
		return Identity{}, fmt.Errorf("wecom http client is not configured")
	}
	token, err := p.corpToken(ctx)
	if err != nil {
		return Identity{}, err
	}
	user, err := p.userFromCode(ctx, token, code)
	if err != nil {
		return Identity{}, err
	}
	name, err := p.userName(ctx, token, user)
	if err != nil {
		return Identity{}, err
	}
	return Identity{UserID: user, Name: name}, nil
}

func (p Production) corpToken(ctx context.Context) (string, error) {
	u := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", url.QueryEscape(p.CorpID), url.QueryEscape(p.Secret))
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

func (p Production) userFromCode(ctx context.Context, token, code string) (string, error) {
	u := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s", url.QueryEscape(token), url.QueryEscape(code))
	body, err := p.HTTP.Get(ctx, u)
	if err != nil {
		return "", fmt.Errorf("wecom userinfo: %w", err)
	}
	var out userResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("wecom userinfo decode: %w", err)
	}
	if out.ErrCode != 0 || out.UserID == "" {
		return "", fmt.Errorf("wecom userinfo error")
	}
	return out.UserID, nil
}

func (p Production) userName(ctx context.Context, token, userid string) (string, error) {
	u := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s", url.QueryEscape(token), url.QueryEscape(userid))
	body, err := p.HTTP.Get(ctx, u)
	if err != nil {
		return "", fmt.Errorf("wecom user get: %w", err)
	}
	var out userGetResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("wecom user get decode: %w", err)
	}
	if out.ErrCode != 0 {
		return "", fmt.Errorf("wecom user get error")
	}
	if out.Name == "" {
		return userid, nil
	}
	return out.Name, nil
}
