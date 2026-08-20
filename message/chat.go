package message

import (
	"context"
	"net/url"
)

// AppChat is an application group chat.
type AppChat struct {
	ChatID   string   `json:"chatid,omitempty"`
	Name     string   `json:"name,omitempty"`
	Owner    string   `json:"owner,omitempty"`
	UserList []string `json:"userlist,omitempty"`
}

// CreateAppChat creates an application group chat and returns chatid.
func (a *API) CreateAppChat(ctx context.Context, chat AppChat) (string, error) {
	var out struct {
		ChatID string `json:"chatid"`
	}
	if err := a.App.PostJSON(ctx, "/cgi-bin/appchat/create", chat, &out); err != nil {
		return "", err
	}
	return out.ChatID, nil
}

// UpdateAppChat updates an application group chat.
func (a *API) UpdateAppChat(ctx context.Context, chat AppChat) error {
	return a.App.PostJSON(ctx, "/cgi-bin/appchat/update", chat, nil)
}

// GetAppChat returns an application group chat.
func (a *API) GetAppChat(ctx context.Context, chatID string) (AppChat, error) {
	var out struct {
		ChatInfo AppChat `json:"chat_info"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/appchat/get", url.Values{"chatid": {chatID}}, &out); err != nil {
		return AppChat{}, err
	}
	return out.ChatInfo, nil
}

// SendAppChat sends a message into an application group chat.
func (a *API) SendAppChat(ctx context.Context, payload any) error {
	return a.App.PostJSON(ctx, "/cgi-bin/appchat/send", payload, nil)
}
