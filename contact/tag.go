package contact

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
func (a *API) CreateTag(ctx context.Context, name string, id int) (int, error) {
	body := map[string]any{"tagname": name}
	if id != 0 {
		body["tagid"] = id
	}
	var out struct {
		TagID int `json:"tagid"`
	}
	if err := a.App.PostJSON(ctx, "/cgi-bin/tag/create", body, &out); err != nil {
		return 0, err
	}
	return out.TagID, nil
}

// UpdateTag renames a tag.
func (a *API) UpdateTag(ctx context.Context, id int, name string) error {
	return a.App.PostJSON(ctx, "/cgi-bin/tag/update", map[string]any{"tagid": id, "tagname": name}, nil)
}

// DeleteTag deletes a tag.
func (a *API) DeleteTag(ctx context.Context, id int) error {
	return a.App.GetJSON(ctx, "/cgi-bin/tag/delete", url.Values{"tagid": {itoa(id)}}, nil)
}

// TagMembers is the result of tag/get.
type TagMembers struct {
	Name      string
	UserList  []TagUser
	PartyList []int
}

// TagUser is a user in a tag.
type TagUser struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

// GetTagMembers lists users and departments in a tag.
func (a *API) GetTagMembers(ctx context.Context, id int) (TagMembers, error) {
	var out struct {
		TagName   string    `json:"tagname"`
		UserList  []TagUser `json:"userlist"`
		PartyList []int     `json:"partylist"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/tag/get", url.Values{"tagid": {itoa(id)}}, &out); err != nil {
		return TagMembers{}, err
	}
	return TagMembers{Name: out.TagName, UserList: out.UserList, PartyList: out.PartyList}, nil
}

// AddTagMembers adds users and/or departments to a tag.
func (a *API) AddTagMembers(ctx context.Context, id int, userlist []string, partylist []int) error {
	return a.App.PostJSON(ctx, "/cgi-bin/tag/addtagusers", tagMemberBody(id, userlist, partylist), nil)
}

// DeleteTagMembers removes users and/or departments from a tag.
func (a *API) DeleteTagMembers(ctx context.Context, id int, userlist []string, partylist []int) error {
	return a.App.PostJSON(ctx, "/cgi-bin/tag/deltagusers", tagMemberBody(id, userlist, partylist), nil)
}

// ListTags returns all tags.
func (a *API) ListTags(ctx context.Context) ([]Tag, error) {
	var out struct {
		TagList []Tag `json:"taglist"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/tag/list", nil, &out); err != nil {
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
