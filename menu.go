package wecom

import (
	"context"
	"net/url"
)

// Menu is the application custom menu.
type Menu struct {
	Button []MenuButton `json:"button"`
}

// MenuButton is a top-level or sub menu item.
type MenuButton struct {
	Type      string       `json:"type,omitempty"`
	Name      string       `json:"name"`
	Key       string       `json:"key,omitempty"`
	URL       string       `json:"url,omitempty"`
	AppID     string       `json:"appid,omitempty"`
	PagePath  string       `json:"pagepath,omitempty"`
	SubButton []MenuButton `json:"sub_button,omitempty"`
}

// CreateMenu creates the application menu.
func (c *Client) CreateMenu(ctx context.Context, menu Menu) error {
	q := url.Values{}
	q.Set("agentid", itoa(c.AgentID))
	return c.postQuery(ctx, "/cgi-bin/menu/create", q, menu, nil)
}

// GetMenu returns the application menu.
func (c *Client) GetMenu(ctx context.Context) (Menu, error) {
	q := url.Values{}
	q.Set("agentid", itoa(c.AgentID))
	var out struct {
		apiMeta
		Menu
	}
	if err := c.get(ctx, "/cgi-bin/menu/get", q, &out); err != nil {
		return Menu{}, err
	}
	return out.Menu, nil
}

// DeleteMenu deletes the application menu.
func (c *Client) DeleteMenu(ctx context.Context) error {
	q := url.Values{}
	q.Set("agentid", itoa(c.AgentID))
	return c.get(ctx, "/cgi-bin/menu/delete", q, nil)
}
