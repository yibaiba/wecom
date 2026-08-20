package wecom

import (
	"context"
	"fmt"
	"net/url"
)

// Identity is the verified WeCom member.
type Identity struct {
	UserID   string
	Name     string
	Email    string
	Avatar   string
	Mobile   string
	Alias    string
	Position string
	Gender   string
}

// Exchanger converts a WeCom authorization code into an identity.
type Exchanger interface {
	AuthURL(state, redirectURI string) string
	Exchange(ctx context.Context, code string) (Identity, error)
	Mode() string
}

// Sandbox is the explicit local adapter. It must be selected in configuration.
type Sandbox struct {
	Users       map[string]Identity
	CallbackURL string
}

// Mode returns sandbox.
func (s Sandbox) Mode() string { return "sandbox" }

// AuthURL returns the local picker/callback URL carrying the host-generated state.
func (s Sandbox) AuthURL(state, redirectURI string) string {
	u, _ := url.Parse(s.CallbackURL)
	q := u.Query()
	q.Set("state", state)
	if redirectURI != "" {
		q.Set("redirect_uri", redirectURI)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// Exchange treats the sandbox code as a configured userid.
func (s Sandbox) Exchange(_ context.Context, code string) (Identity, error) {
	id, ok := s.Users[code]
	if !ok {
		return Identity{}, fmt.Errorf("unknown sandbox userid")
	}
	return id, nil
}

// Production talks to the WeCom OAuth API.
type Production struct {
	CorpID      string
	AgentID     int
	Secret      string
	HTTP        Doer
	RedirectURI string
}

// Doer is the injected HTTP client.
type Doer interface {
	Get(ctx context.Context, rawURL string) ([]byte, error)
}

// Mode returns production.
func (p Production) Mode() string { return "production" }

// AuthURL is the official WeCom QR connect URL.
func (p Production) AuthURL(state, redirectURI string) string {
	redir := redirectURI
	if redir == "" {
		redir = p.RedirectURI
	}
	u, _ := url.Parse("https://open.work.weixin.qq.com/wwopen/sso/qrConnect")
	q := u.Query()
	q.Set("appid", p.CorpID)
	q.Set("agentid", fmt.Sprintf("%d", p.AgentID))
	q.Set("redirect_uri", redir)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String()
}

// WxWorkAuthURL is the in-WeCom browser OAuth2 authorize URL.
func WxWorkAuthURL(corpID string, agentID int, redirectURI, state string) string {
	return wxWorkAuthorizeURL(corpID, agentID, redirectURI, state, "snsapi_base")
}

// WxWorkPrivateInfoURL is the phone-scan OAuth2 URL for unknown members.
func WxWorkPrivateInfoURL(corpID string, agentID int, redirectURI, state string) string {
	return wxWorkAuthorizeURL(corpID, agentID, redirectURI, state, "snsapi_privateinfo")
}

func wxWorkAuthorizeURL(corpID string, agentID int, redirectURI, state, scope string) string {
	u, _ := url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")
	q := u.Query()
	q.Set("appid", corpID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", scope)
	q.Set("state", state)
	q.Set("agentid", fmt.Sprintf("%d", agentID))
	u.RawQuery = q.Encode()
	u.Fragment = "wechat_redirect"
	return u.String()
}
