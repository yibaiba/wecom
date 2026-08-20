package agent

import (
	"context"
	"net/url"
	"strconv"
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

func (a *API) agentQuery() url.Values {
	return url.Values{"agentid": {strconv.Itoa(a.App.AgentID)}}
}

// CreateMenu creates the application menu.
func (a *API) CreateMenu(ctx context.Context, menu Menu) error {
	return a.App.PostQuery(ctx, "/cgi-bin/menu/create", a.agentQuery(), menu, nil)
}

// GetMenu returns the application menu.
func (a *API) GetMenu(ctx context.Context) (Menu, error) {
	var out Menu
	if err := a.App.GetJSON(ctx, "/cgi-bin/menu/get", a.agentQuery(), &out); err != nil {
		return Menu{}, err
	}
	return out, nil
}

// DeleteMenu deletes the application menu.
func (a *API) DeleteMenu(ctx context.Context) error {
	return a.App.GetJSON(ctx, "/cgi-bin/menu/delete", a.agentQuery(), nil)
}
