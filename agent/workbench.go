package agent

import "context"

// SetWorkbenchTemplate sets the workbench home layout. payload is the official JSON body.
func (a *API) SetWorkbenchTemplate(ctx context.Context, payload any) error {
	return a.App.PostJSON(ctx, "/cgi-bin/agent/set_workbench_template", payload, nil)
}

// GetWorkbenchTemplate returns the workbench home layout JSON.
func (a *API) GetWorkbenchTemplate(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := a.App.GetJSON(ctx, "/cgi-bin/agent/get_workbench_template", a.agentQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetWorkbenchData sets a member's workbench data. payload is the official JSON body.
func (a *API) SetWorkbenchData(ctx context.Context, payload any) error {
	return a.App.PostJSON(ctx, "/cgi-bin/agent/set_workbench_data", payload, nil)
}
