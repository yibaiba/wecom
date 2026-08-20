package message

import "github.com/yibaiba/wecom"

// API sends application messages and manages app group chats.
type API struct {
	App *wecom.Client
}

// New wraps a shared WeCom client.
func New(app *wecom.Client) *API {
	return &API{App: app}
}
