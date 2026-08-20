package contact

import (
	"context"
	"net/url"
)

// Department is a contacts department.
type Department struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	NameEN           string   `json:"name_en,omitempty"`
	ParentID         int      `json:"parentid"`
	Order            int      `json:"order"`
	DepartmentLeader []string `json:"department_leader,omitempty"`
}

// ListDepartments returns departments visible to the app. id 0 means the full tree.
func (a *API) ListDepartments(ctx context.Context, id int) ([]Department, error) {
	q := url.Values{}
	if id != 0 {
		q.Set("id", itoa(id))
	}
	var out struct {
		Department []Department `json:"department"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/department/list", q, &out); err != nil {
		return nil, err
	}
	return out.Department, nil
}

// ListDepartmentIDs returns child department IDs (simplelist).
func (a *API) ListDepartmentIDs(ctx context.Context, id int) ([]Department, error) {
	q := url.Values{}
	if id != 0 {
		q.Set("id", itoa(id))
	}
	var out struct {
		DepartmentID []Department `json:"department_id"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/department/simplelist", q, &out); err != nil {
		return nil, err
	}
	return out.DepartmentID, nil
}

// GetDepartment returns one department.
func (a *API) GetDepartment(ctx context.Context, id int) (Department, error) {
	q := url.Values{}
	q.Set("id", itoa(id))
	var out struct {
		Department Department `json:"department"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/department/get", q, &out); err != nil {
		return Department{}, err
	}
	return out.Department, nil
}

// CreateDepartment creates a department and returns its id.
func (a *API) CreateDepartment(ctx context.Context, d Department) (int, error) {
	var out struct {
		ID int `json:"id"`
	}
	if err := a.App.PostJSON(ctx, "/cgi-bin/department/create", d, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

// DepartmentPatch is a partial department update. Nil pointer fields are omitted
// so parentid/order 0 is not posted unless the caller sets the pointer.
type DepartmentPatch struct {
	ID               int      `json:"id"`
	Name             string   `json:"name,omitempty"`
	NameEN           string   `json:"name_en,omitempty"`
	ParentID         *int     `json:"parentid,omitempty"`
	Order            *int     `json:"order,omitempty"`
	DepartmentLeader []string `json:"department_leader,omitempty"`
}

// UpdateDepartment updates a department. Unset ParentID/Order are omitted.
func (a *API) UpdateDepartment(ctx context.Context, d DepartmentPatch) error {
	return a.App.PostJSON(ctx, "/cgi-bin/department/update", d, nil)
}

// DeleteDepartment deletes a department.
func (a *API) DeleteDepartment(ctx context.Context, id int) error {
	q := url.Values{}
	q.Set("id", itoa(id))
	return a.App.GetJSON(ctx, "/cgi-bin/department/delete", q, nil)
}
