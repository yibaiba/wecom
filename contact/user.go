package contact

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/yibaiba/wecom"
)

// UserBrief is one row from user/simplelist.
type UserBrief struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	OpenUserID string `json:"open_userid"`
}

// ListDepartmentUsers lists members in a department (user/simplelist).
func (a *API) ListDepartmentUsers(ctx context.Context, departmentID int) ([]UserBrief, error) {
	q := url.Values{}
	q.Set("department_id", itoa(departmentID))
	var out struct {
		UserList []UserBrief `json:"userlist"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/user/simplelist", q, &out); err != nil {
		return nil, err
	}
	return out.UserList, nil
}

// ListDepartmentUserDetails lists full members in a department (user/list).
func (a *API) ListDepartmentUserDetails(ctx context.Context, departmentID int) ([]wecom.Identity, error) {
	q := url.Values{}
	q.Set("department_id", itoa(departmentID))
	var raw json.RawMessage
	if err := a.App.GetJSON(ctx, "/cgi-bin/user/list", q, &raw); err != nil {
		return nil, err
	}
	return wecom.ParseDirectoryUsers(raw)
}

// GetUser reads one member from the directory.
func (a *API) GetUser(ctx context.Context, userid string) (wecom.Identity, error) {
	q := url.Values{}
	q.Set("userid", userid)
	var raw json.RawMessage
	if err := a.App.GetJSON(ctx, "/cgi-bin/user/get", q, &raw); err != nil {
		return wecom.Identity{}, err
	}
	return wecom.ParseDirectoryUser(userid, raw)
}

// MemberInput is the body for user/create and user/update.
type MemberInput struct {
	UserID           string                 `json:"userid"`
	Name             string                 `json:"name,omitempty"`
	Alias            string                 `json:"alias,omitempty"`
	Mobile           string                 `json:"mobile,omitempty"`
	Department       []int                  `json:"department,omitempty"`
	Order            []int                  `json:"order,omitempty"`
	Position         string                 `json:"position,omitempty"`
	Gender           string                 `json:"gender,omitempty"`
	Email            string                 `json:"email,omitempty"`
	BizMail          string                 `json:"biz_mail,omitempty"`
	IsLeaderInDept   []int                  `json:"is_leader_in_dept,omitempty"`
	DirectLeader     []string               `json:"direct_leader,omitempty"`
	Enable           *int                   `json:"enable,omitempty"`
	AvatarMediaID    string                 `json:"avatar_mediaid,omitempty"`
	Telephone        string                 `json:"telephone,omitempty"`
	Address          string                 `json:"address,omitempty"`
	MainDepartment   int                    `json:"main_department,omitempty"`
	ToInvite         *bool                  `json:"to_invite,omitempty"`
	ExternalPosition string                 `json:"external_position,omitempty"`
	ExtAttr          *wecom.ExtAttrGroup    `json:"extattr,omitempty"`
	ExternalProfile  *wecom.ExternalProfile `json:"external_profile,omitempty"`
}

// CreateUser creates a directory member.
func (a *API) CreateUser(ctx context.Context, in MemberInput) error {
	return a.App.PostJSON(ctx, "/cgi-bin/user/create", in, nil)
}

// UpdateUser updates a directory member.
func (a *API) UpdateUser(ctx context.Context, in MemberInput) error {
	return a.App.PostJSON(ctx, "/cgi-bin/user/update", in, nil)
}

// DeleteUser deletes a directory member.
func (a *API) DeleteUser(ctx context.Context, userid string) error {
	q := url.Values{}
	q.Set("userid", userid)
	return a.App.GetJSON(ctx, "/cgi-bin/user/delete", q, nil)
}

// BatchDeleteUsers deletes members.
func (a *API) BatchDeleteUsers(ctx context.Context, userids []string) error {
	return a.App.PostJSON(ctx, "/cgi-bin/user/batchdelete", map[string][]string{"useridlist": userids}, nil)
}
