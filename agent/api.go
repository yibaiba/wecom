package agent

import "github.com/yibaiba/wecom"

// API is application settings, menu, and workbench.
type API struct {
	App *wecom.Client
}

// New wraps a shared WeCom client.
func New(app *wecom.Client) *API {
	return &API{App: app}
}
