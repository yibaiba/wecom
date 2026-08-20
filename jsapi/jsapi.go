package jsapi

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"

	"github.com/yibaiba/wecom"
)

// API fetches JS-SDK tickets.
type API struct {
	App *wecom.Client
}

// New wraps a shared WeCom client.
func New(app *wecom.Client) *API {
	return &API{App: app}
}

type ticketResp struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// Ticket returns the corp jsapi_ticket for ww.register / getConfigSignature.
func (a *API) Ticket(ctx context.Context) (string, error) {
	var out ticketResp
	if err := a.App.GetJSON(ctx, "/cgi-bin/get_jsapi_ticket", nil, &out); err != nil {
		return "", err
	}
	return out.Ticket, nil
}

// AgentTicket returns the agent_config ticket for getAgentConfigSignature.
func (a *API) AgentTicket(ctx context.Context) (string, error) {
	q := url.Values{}
	q.Set("type", "agent_config")
	var out ticketResp
	if err := a.App.GetJSON(ctx, "/cgi-bin/ticket/get", q, &out); err != nil {
		return "", err
	}
	return out.Ticket, nil
}

// Signature is the official SHA1 signature string.
func Signature(ticket, nonce, rawURL string, timestamp int64) string {
	plain := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, rawURL)
	sum := sha1.Sum([]byte(plain))
	return fmt.Sprintf("%x", sum)
}
