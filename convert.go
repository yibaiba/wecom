package wecom

import (
	"context"
	"fmt"
	"net/url"
)

func fmtInt(n int) string { return fmt.Sprintf("%d", n) }

// UserIDToOpenID converts a userid to an openid.
func (c *Client) UserIDToOpenID(ctx context.Context, userid string) (string, error) {
	var out struct {
		apiMeta
		OpenID string `json:"openid"`
	}
	if err := c.post(ctx, "/cgi-bin/user/convert_to_openid", map[string]string{"userid": userid}, &out); err != nil {
		return "", err
	}
	return out.OpenID, nil
}

// OpenIDToUserID converts an openid to a userid.
func (c *Client) OpenIDToUserID(ctx context.Context, openid string) (string, error) {
	var out struct {
		apiMeta
		UserID string `json:"userid"`
	}
	if err := c.post(ctx, "/cgi-bin/user/convert_to_userid", map[string]string{"openid": openid}, &out); err != nil {
		return "", err
	}
	return out.UserID, nil
}

// AuthSuccess marks login secondary verification as complete.
func (c *Client) AuthSuccess(ctx context.Context, userid string) error {
	q := url.Values{}
	q.Set("userid", userid)
	return c.get(ctx, "/cgi-bin/user/authsucc", q, nil)
}

// InviteUsers sends an invitation to join WeCom.
func (c *Client) InviteUsers(ctx context.Context, user, party, tag []string) error {
	body := map[string]any{}
	if len(user) > 0 {
		body["user"] = user
	}
	if len(party) > 0 {
		body["party"] = party
	}
	if len(tag) > 0 {
		body["tag"] = tag
	}
	return c.post(ctx, "/cgi-bin/batch/invite", body, nil)
}

// UserIDByMobile looks up a userid by mobile.
func (c *Client) UserIDByMobile(ctx context.Context, mobile string) (string, error) {
	var out struct {
		apiMeta
		UserID string `json:"userid"`
	}
	if err := c.post(ctx, "/cgi-bin/user/getuserid", map[string]string{"mobile": mobile}, &out); err != nil {
		return "", err
	}
	return out.UserID, nil
}

// UserIDByEmail looks up a userid by email. emailType 1 personal, 2 biz mail.
func (c *Client) UserIDByEmail(ctx context.Context, email string, emailType int) (string, error) {
	if emailType == 0 {
		emailType = 1
	}
	var out struct {
		apiMeta
		UserID string `json:"userid"`
	}
	body := map[string]any{"email": email, "email_type": emailType}
	if err := c.post(ctx, "/cgi-bin/user/get_userid_by_email", body, &out); err != nil {
		return "", err
	}
	return out.UserID, nil
}

// UserIDPage is one page from user/list_id.
type UserIDPage struct {
	NextCursor string
	DeptUsers  []DeptUser
}

// DeptUser is a userid plus department from list_id.
type DeptUser struct {
	UserID     string `json:"userid"`
	Department int    `json:"department"`
}

// ListUserIDs pages through member IDs.
func (c *Client) ListUserIDs(ctx context.Context, cursor string, limit int) (UserIDPage, error) {
	if limit == 0 {
		limit = 1000
	}
	var out struct {
		apiMeta
		NextCursor string     `json:"next_cursor"`
		DeptUser   []DeptUser `json:"dept_user"`
	}
	body := map[string]any{"limit": limit}
	if cursor != "" {
		body["cursor"] = cursor
	}
	if err := c.post(ctx, "/cgi-bin/user/list_id", body, &out); err != nil {
		return UserIDPage{}, err
	}
	return UserIDPage{NextCursor: out.NextCursor, DeptUsers: out.DeptUser}, nil
}

// JoinQRCode returns the join-enterprise QR URL. sizeType 1-4.
func (c *Client) JoinQRCode(ctx context.Context, sizeType int) (string, error) {
	q := url.Values{}
	if sizeType != 0 {
		q.Set("size_type", itoa(sizeType))
	}
	var out struct {
		apiMeta
		JoinQRCode string `json:"join_qrcode"`
	}
	if err := c.get(ctx, "/cgi-bin/corp/get_join_qrcode", q, &out); err != nil {
		return "", err
	}
	return out.JoinQRCode, nil
}
