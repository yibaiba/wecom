package wecom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"
)

var apiBase = "https://qyapi.weixin.qq.com"

// Client is a WeCom self-built application. Hosts that only need login can
// keep using Production; Client covers the rest of the app APIs.
type Client struct {
	CorpID  string
	AgentID int
	Secret  string
	HTTP    Doer

	now      func() time.Time
	mu       sync.Mutex
	token    string
	tokenExp time.Time
}

// App returns a Client with the same corp credentials as Production.
func (p Production) App() *Client {
	return &Client{CorpID: p.CorpID, AgentID: p.AgentID, Secret: p.Secret, HTTP: p.HTTP}
}

func (c *Client) clock() time.Time {
	if c.now != nil {
		return c.now()
	}
	return time.Now()
}

func (c *Client) doer() Doer {
	if c.HTTP != nil {
		return c.HTTP
	}
	return HTTPDoer{}
}

// Token returns a cached access_token.
func (c *Client) Token(ctx context.Context) (string, error) {
	now := c.clock()
	c.mu.Lock()
	if c.token != "" && now.Before(c.tokenExp) {
		tok := c.token
		c.mu.Unlock()
		return tok, nil
	}
	c.mu.Unlock()
	tok, exp, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.token, c.tokenExp = tok, exp
	c.mu.Unlock()
	return tok, nil
}

func (c *Client) invalidateToken() {
	c.mu.Lock()
	c.token, c.tokenExp = "", time.Time{}
	c.mu.Unlock()
}

func (c *Client) fetchToken(ctx context.Context) (string, time.Time, error) {
	q := url.Values{}
	q.Set("corpid", c.CorpID)
	q.Set("corpsecret", c.Secret)
	body, err := c.doer().Get(ctx, apiBase+"/cgi-bin/gettoken?"+q.Encode())
	if err != nil {
		return "", time.Time{}, fmt.Errorf("wecom token: %w", err)
	}
	var out tokenResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", time.Time{}, fmt.Errorf("wecom token decode: %w", err)
	}
	if out.ErrCode != 0 || out.AccessToken == "" {
		return "", time.Time{}, Error{Code: out.ErrCode, Msg: firstNonEmpty(out.ErrMsg, "token error")}
	}
	sec := out.ExpiresIn
	if sec <= 0 {
		sec = 7200
	}
	if sec > 400 {
		sec -= 400
	}
	return out.AccessToken, c.clock().Add(time.Duration(sec) * time.Second), nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, dest any) error {
	return c.call(ctx, false, path, query, nil, dest)
}

func (c *Client) post(ctx context.Context, path string, payload any, dest any) error {
	return c.postQuery(ctx, path, nil, payload, dest)
}

func (c *Client) postQuery(ctx context.Context, path string, query url.Values, payload any, dest any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.call(ctx, true, path, query, body, dest)
}

func (c *Client) call(ctx context.Context, post bool, path string, query url.Values, payload []byte, dest any) error {
	err := c.doCall(ctx, post, path, query, payload, dest)
	if !isTokenErr(err) {
		return err
	}
	c.invalidateToken()
	return c.doCall(ctx, post, path, query, payload, dest)
}

func (c *Client) doCall(ctx context.Context, post bool, path string, query url.Values, payload []byte, dest any) error {
	tok, err := c.Token(ctx)
	if err != nil {
		return err
	}
	if query == nil {
		query = url.Values{}
	}
	query.Set("access_token", tok)
	raw := apiBase + path + "?" + query.Encode()
	var body []byte
	if post {
		body, err = c.doer().Post(ctx, raw, payload)
	} else {
		body, err = c.doer().Get(ctx, raw)
	}
	if err != nil {
		return err
	}
	if err := decodeAPI(body, dest); err != nil {
		return err
	}
	return nil
}

type apiMeta struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func decodeAPI(body []byte, dest any) error {
	var meta apiMeta
	if err := json.Unmarshal(body, &meta); err != nil {
		return fmt.Errorf("wecom decode: %w", err)
	}
	if meta.ErrCode != 0 {
		return Error{Code: meta.ErrCode, Msg: meta.ErrMsg}
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("wecom decode: %w", err)
	}
	return nil
}
