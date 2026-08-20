package wecom

import "context"

// AppChat is an application group chat.
type AppChat struct {
	ChatID   string   `json:"chatid,omitempty"`
	Name     string   `json:"name,omitempty"`
	Owner    string   `json:"owner,omitempty"`
	UserList []string `json:"userlist,omitempty"`
}

// CreateAppChat creates an application group chat and returns chatid.
func (c *Client) CreateAppChat(ctx context.Context, chat AppChat) (string, error) {
	var out struct {
		apiMeta
		ChatID string `json:"chatid"`
	}
	if err := c.post(ctx, "/cgi-bin/appchat/create", chat, &out); err != nil {
		return "", err
	}
	return out.ChatID, nil
}

// UpdateAppChat updates an application group chat.
func (c *Client) UpdateAppChat(ctx context.Context, chat AppChat) error {
	return c.post(ctx, "/cgi-bin/appchat/update", chat, nil)
}

// GetAppChat returns an application group chat.
func (c *Client) GetAppChat(ctx context.Context, chatID string) (AppChat, error) {
	var out struct {
		apiMeta
		ChatInfo AppChat `json:"chat_info"`
	}
	if err := c.get(ctx, "/cgi-bin/appchat/get", mapValues("chatid", chatID), &out); err != nil {
		return AppChat{}, err
	}
	return out.ChatInfo, nil
}

// SendAppChat sends a message into an application group chat.
func (c *Client) SendAppChat(ctx context.Context, payload any) error {
	return c.post(ctx, "/cgi-bin/appchat/send", payload, nil)
}

// SetWorkbenchTemplate sets the workbench home layout. payload is the official JSON body.
func (c *Client) SetWorkbenchTemplate(ctx context.Context, payload any) error {
	return c.post(ctx, "/cgi-bin/agent/set_workbench_template", payload, nil)
}

// GetWorkbenchTemplate returns the workbench home layout JSON.
func (c *Client) GetWorkbenchTemplate(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := c.get(ctx, "/cgi-bin/agent/get_workbench_template", mapValues("agentid", itoa(c.AgentID)), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetWorkbenchData sets a member's workbench data. payload is the official JSON body.
func (c *Client) SetWorkbenchData(ctx context.Context, payload any) error {
	return c.post(ctx, "/cgi-bin/agent/set_workbench_data", payload, nil)
}
