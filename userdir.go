package wecom

import (
	"context"
	"net/url"
)

// UserBrief is one row from user/simplelist.
type UserBrief struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	OpenUserID string `json:"open_userid"`
}

// ListDepartmentUsers lists members in a department (user/simplelist).
func (c *Client) ListDepartmentUsers(ctx context.Context, departmentID int) ([]UserBrief, error) {
	q := url.Values{}
	q.Set("department_id", itoa(departmentID))
	var out struct {
		apiMeta
		UserList []UserBrief `json:"userlist"`
	}
	if err := c.get(ctx, "/cgi-bin/user/simplelist", q, &out); err != nil {
		return nil, err
	}
	return out.UserList, nil
}

// ListDepartmentUserDetails lists full members in a department (user/list).
func (c *Client) ListDepartmentUserDetails(ctx context.Context, departmentID int) ([]Identity, error) {
	q := url.Values{}
	q.Set("department_id", itoa(departmentID))
	var out struct {
		apiMeta
		UserList []userGetResp `json:"userlist"`
	}
	if err := c.get(ctx, "/cgi-bin/user/list", q, &out); err != nil {
		return nil, err
	}
	users := make([]Identity, 0, len(out.UserList))
	for _, row := range out.UserList {
		users = append(users, identityFromUserGet(row.UserID, row))
	}
	return users, nil
}

// GetUser reads one member from the directory.
func (c *Client) GetUser(ctx context.Context, userid string) (Identity, error) {
	q := url.Values{}
	q.Set("userid", userid)
	var out userGetResp
	if err := c.get(ctx, "/cgi-bin/user/get", q, &out); err != nil {
		return Identity{}, err
	}
	return identityFromUserGet(userid, out), nil
}

// MemberInput is the body for user/create and user/update.
type MemberInput struct {
	UserID           string           `json:"userid"`
	Name             string           `json:"name,omitempty"`
	Alias            string           `json:"alias,omitempty"`
	Mobile           string           `json:"mobile,omitempty"`
	Department       []int            `json:"department,omitempty"`
	Order            []int            `json:"order,omitempty"`
	Position         string           `json:"position,omitempty"`
	Gender           string           `json:"gender,omitempty"`
	Email            string           `json:"email,omitempty"`
	BizMail          string           `json:"biz_mail,omitempty"`
	IsLeaderInDept   []int            `json:"is_leader_in_dept,omitempty"`
	DirectLeader     []string         `json:"direct_leader,omitempty"`
	Enable           *int             `json:"enable,omitempty"`
	AvatarMediaID    string           `json:"avatar_mediaid,omitempty"`
	Telephone        string           `json:"telephone,omitempty"`
	Address          string           `json:"address,omitempty"`
	MainDepartment   int              `json:"main_department,omitempty"`
	ToInvite         *bool            `json:"to_invite,omitempty"`
	ExternalPosition string           `json:"external_position,omitempty"`
	ExtAttr          *ExtAttrGroup    `json:"extattr,omitempty"`
	ExternalProfile  *ExternalProfile `json:"external_profile,omitempty"`
}

// CreateUser creates a directory member.
func (c *Client) CreateUser(ctx context.Context, in MemberInput) error {
	return c.post(ctx, "/cgi-bin/user/create", in, nil)
}

// UpdateUser updates a directory member.
func (c *Client) UpdateUser(ctx context.Context, in MemberInput) error {
	return c.post(ctx, "/cgi-bin/user/update", in, nil)
}

// DeleteUser deletes a directory member.
func (c *Client) DeleteUser(ctx context.Context, userid string) error {
	q := url.Values{}
	q.Set("userid", userid)
	return c.get(ctx, "/cgi-bin/user/delete", q, nil)
}

// BatchDeleteUsers deletes members.
func (c *Client) BatchDeleteUsers(ctx context.Context, userids []string) error {
	return c.post(ctx, "/cgi-bin/user/batchdelete", map[string][]string{"useridlist": userids}, nil)
}

func itoa(n int) string {
	return fmtInt(n)
}
