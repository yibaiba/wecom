package agent

import "context"

// WorkbenchData is a workbench template or per-user payload.
// Type is keydata, image, list, webview, or normal.
type WorkbenchData struct {
	Type    string            `json:"type"`
	KeyData *WorkbenchKeyData `json:"keydata,omitempty"`
	Image   *WorkbenchImage   `json:"image,omitempty"`
	List    *WorkbenchList    `json:"list,omitempty"`
	Webview *WorkbenchWebview `json:"webview,omitempty"`
}

// WorkbenchKeyData is type "keydata" (at most 4 items).
type WorkbenchKeyData struct {
	Items []WorkbenchKeyItem `json:"items"`
}

// WorkbenchKeyItem is one keydata cell.
type WorkbenchKeyItem struct {
	Key      string `json:"key,omitempty"`
	Data     string `json:"data"`
	JumpURL  string `json:"jump_url,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
}

// WorkbenchImage is type "image".
type WorkbenchImage struct {
	URL      string `json:"url"`
	JumpURL  string `json:"jump_url,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
}

// WorkbenchList is type "list" (at most 3 items).
type WorkbenchList struct {
	Items []WorkbenchListItem `json:"items"`
}

// WorkbenchListItem is one list row.
type WorkbenchListItem struct {
	Title    string `json:"title"`
	JumpURL  string `json:"jump_url,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
}

// WorkbenchWebview is type "webview". Height is single_row or double_row.
type WorkbenchWebview struct {
	URL                string `json:"url"`
	JumpURL            string `json:"jump_url,omitempty"`
	PagePath           string `json:"pagepath,omitempty"`
	Height             string `json:"height,omitempty"`
	HideTitle          bool   `json:"hide_title,omitempty"`
	EnableWebviewClick bool   `json:"enable_webview_click,omitempty"`
}

func (a *API) agentBody() map[string]any {
	return map[string]any{"agentid": a.App.AgentID}
}

// SetWorkbenchTemplate sets the workbench home layout. payload is the official JSON body.
func (a *API) SetWorkbenchTemplate(ctx context.Context, payload any) error {
	return a.App.PostJSON(ctx, "/cgi-bin/agent/set_workbench_template", payload, nil)
}

// GetWorkbenchTemplate returns the workbench home layout JSON.
func (a *API) GetWorkbenchTemplate(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := a.App.PostJSON(ctx, "/cgi-bin/agent/get_workbench_template", a.agentBody(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetWorkbenchData sets a member's workbench data. payload is the official JSON body.
func (a *API) SetWorkbenchData(ctx context.Context, payload any) error {
	return a.App.PostJSON(ctx, "/cgi-bin/agent/set_workbench_data", payload, nil)
}

// BatchSetWorkbenchData sets the same workbench data for many members.
func (a *API) BatchSetWorkbenchData(ctx context.Context, userIDs []string, data WorkbenchData) error {
	body := map[string]any{
		"agentid":     a.App.AgentID,
		"userid_list": userIDs,
		"data":        data,
	}
	return a.App.PostJSON(ctx, "/cgi-bin/agent/batch_set_workbench_data", body, nil)
}

// GetWorkbenchData returns one member's workbench data.
func (a *API) GetWorkbenchData(ctx context.Context, userid string) (WorkbenchData, error) {
	body := map[string]any{
		"agentid": a.App.AgentID,
		"userid":  userid,
	}
	var out struct {
		Data WorkbenchData `json:"data"`
	}
	if err := a.App.PostJSON(ctx, "/cgi-bin/agent/get_workbench_data", body, &out); err != nil {
		return WorkbenchData{}, err
	}
	return out.Data, nil
}
