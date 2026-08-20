package contact

import (
	"strconv"

	"github.com/yibaiba/wecom"
)

// API is the address-book surface: members, departments, tags.
type API struct {
	App *wecom.Client
}

// New wraps a shared WeCom client.
func New(app *wecom.Client) *API {
	return &API{App: app}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
