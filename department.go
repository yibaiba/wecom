package wecom

import (
	"context"
	"fmt"
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
func (c *Client) ListDepartments(ctx context.Context, id int) ([]Department, error) {
	q := url.Values{}
	if id != 0 {
		q.Set("id", fmt.Sprintf("%d", id))
	}
	var out struct {
		apiMeta
		Department []Department `json:"department"`
	}
	if err := c.get(ctx, "/cgi-bin/department/list", q, &out); err != nil {
		return nil, err
	}
	return out.Department, nil
}

// ListDepartmentIDs returns child department IDs (simplelist).
func (c *Client) ListDepartmentIDs(ctx context.Context, id int) ([]Department, error) {
	q := url.Values{}
	if id != 0 {
		q.Set("id", fmt.Sprintf("%d", id))
	}
	var out struct {
		apiMeta
		DepartmentID []Department `json:"department_id"`
	}
	if err := c.get(ctx, "/cgi-bin/department/simplelist", q, &out); err != nil {
		return nil, err
	}
	return out.DepartmentID, nil
}

// GetDepartment returns one department.
func (c *Client) GetDepartment(ctx context.Context, id int) (Department, error) {
	q := url.Values{}
	q.Set("id", fmt.Sprintf("%d", id))
	var out struct {
		apiMeta
		Department Department `json:"department"`
	}
	if err := c.get(ctx, "/cgi-bin/department/get", q, &out); err != nil {
		return Department{}, err
	}
	return out.Department, nil
}

// CreateDepartment creates a department and returns its id.
func (c *Client) CreateDepartment(ctx context.Context, d Department) (int, error) {
	var out struct {
		apiMeta
		ID int `json:"id"`
	}
	if err := c.post(ctx, "/cgi-bin/department/create", d, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

// UpdateDepartment updates a department.
func (c *Client) UpdateDepartment(ctx context.Context, d Department) error {
	return c.post(ctx, "/cgi-bin/department/update", d, nil)
}

// DeleteDepartment deletes a department.
func (c *Client) DeleteDepartment(ctx context.Context, id int) error {
	q := url.Values{}
	q.Set("id", fmt.Sprintf("%d", id))
	return c.get(ctx, "/cgi-bin/department/delete", q, nil)
}
