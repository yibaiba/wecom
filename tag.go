package wecom

import (
	"context"
	"net/url"
)

// Tag is a contacts tag.
type Tag struct {
	ID   int    `json:"tagid,omitempty"`
	Name string `json:"tagname"`
}

// CreateTag creates a tag and returns its id.
func (c *Client) CreateTag(ctx context.Context, name string, id int) (int, error) {
	body := map[string]any{"tagname": name}
	if id != 0 {
		body["tagid"] = id
	}
	var out struct {
		apiMeta
		TagID int `json:"tagid"`
	}
	if err := c.post(ctx, "/cgi-bin/tag/create", body, &out); err != nil {
		return 0, err
	}
	return out.TagID, nil
}

// UpdateTag renames a tag.
func (c *Client) UpdateTag(ctx context.Context, id int, name string) error {
	return c.post(ctx, "/cgi-bin/tag/update", map[string]any{"tagid": id, "tagname": name}, nil)
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.get(ctx, "/cgi-bin/tag/delete", mapValues("tagid", itoa(id)), nil)
}

// TagMembers is the result of tag/get.
type TagMembers struct {
	Name      string `json:"-"`
	UserList  []TagUser
	PartyList []int
}

// TagUser is a user in a tag.
type TagUser struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

// GetTagMembers lists users and departments in a tag.
func (c *Client) GetTagMembers(ctx context.Context, id int) (TagMembers, error) {
	var out struct {
		apiMeta
		TagName   string    `json:"tagname"`
		UserList  []TagUser `json:"userlist"`
		PartyList []int     `json:"partylist"`
	}
	if err := c.get(ctx, "/cgi-bin/tag/get", mapValues("tagid", itoa(id)), &out); err != nil {
		return TagMembers{}, err
	}
	return TagMembers{Name: out.TagName, UserList: out.UserList, PartyList: out.PartyList}, nil
}

// AddTagMembers adds users and/or departments to a tag.
func (c *Client) AddTagMembers(ctx context.Context, id int, userlist []string, partylist []int) error {
	return c.post(ctx, "/cgi-bin/tag/addtagusers", tagMemberBody(id, userlist, partylist), nil)
}

// DeleteTagMembers removes users and/or departments from a tag.
func (c *Client) DeleteTagMembers(ctx context.Context, id int, userlist []string, partylist []int) error {
	return c.post(ctx, "/cgi-bin/tag/deltagusers", tagMemberBody(id, userlist, partylist), nil)
}

// ListTags returns all tags.
func (c *Client) ListTags(ctx context.Context) ([]Tag, error) {
	var out struct {
		apiMeta
		TagList []Tag `json:"taglist"`
	}
	if err := c.get(ctx, "/cgi-bin/tag/list", nil, &out); err != nil {
		return nil, err
	}
	return out.TagList, nil
}

func tagMemberBody(id int, userlist []string, partylist []int) map[string]any {
	body := map[string]any{"tagid": id}
	if len(userlist) > 0 {
		body["userlist"] = userlist
	}
	if len(partylist) > 0 {
		body["partylist"] = partylist
	}
	return body
}

func mapValues(k, v string) url.Values {
	return url.Values{k: {v}}
}
