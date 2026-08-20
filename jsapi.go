package wecom

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
)

type ticketResp struct {
	apiMeta
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// JSAPITicket returns the corp jsapi_ticket for ww.register / getConfigSignature.
func (c *Client) JSAPITicket(ctx context.Context) (string, error) {
	var out ticketResp
	if err := c.get(ctx, "/cgi-bin/get_jsapi_ticket", nil, &out); err != nil {
		return "", err
	}
	return out.Ticket, nil
}

// AgentTicket returns the agent_config ticket for getAgentConfigSignature.
func (c *Client) AgentTicket(ctx context.Context) (string, error) {
	q := url.Values{}
	q.Set("type", "agent_config")
	var out ticketResp
	if err := c.get(ctx, "/cgi-bin/ticket/get", q, &out); err != nil {
		return "", err
	}
	return out.Ticket, nil
}

// JSAPISignature is the official SHA1 signature string.
func JSAPISignature(ticket, nonce, rawURL string, timestamp int64) string {
	plain := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, rawURL)
	sum := sha1.Sum([]byte(plain))
	return fmt.Sprintf("%x", sum)
}
