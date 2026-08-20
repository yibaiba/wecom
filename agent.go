package wecom

import (
	"context"
	"fmt"
	"net/url"
)

// Agent is a self-built application returned by agent/get.
type Agent struct {
	AgentID                 int      `json:"agentid"`
	Name                    string   `json:"name"`
	SquareLogoURL           string   `json:"square_logo_url"`
	Description             string   `json:"description"`
	Close                   int      `json:"close"`
	RedirectDomain          string   `json:"redirect_domain"`
	ReportLocationFlag      int      `json:"report_location_flag"`
	IsReportEnter           int      `json:"isreportenter"`
	HomeURL                 string   `json:"home_url"`
	CustomizedPublishStatus int      `json:"customized_publish_status"`
	AllowUserIDs            []string `json:"-"`
	AllowPartyIDs           []int    `json:"-"`
	AllowTagIDs             []int    `json:"-"`
}

type agentGetResp struct {
	apiMeta
	AgentID                 int    `json:"agentid"`
	Name                    string `json:"name"`
	SquareLogoURL           string `json:"square_logo_url"`
	Description             string `json:"description"`
	Close                   int    `json:"close"`
	RedirectDomain          string `json:"redirect_domain"`
	ReportLocationFlag      int    `json:"report_location_flag"`
	IsReportEnter           int    `json:"isreportenter"`
	HomeURL                 string `json:"home_url"`
	CustomizedPublishStatus int    `json:"customized_publish_status"`
	AllowUserinfos          struct {
		User []struct {
			UserID string `json:"userid"`
		} `json:"user"`
	} `json:"allow_userinfos"`
	AllowPartys struct {
		PartyID []int `json:"partyid"`
	} `json:"allow_partys"`
	AllowTags struct {
		TagID []int `json:"tagid"`
	} `json:"allow_tags"`
}

// GetAgent returns the current application (or agentID if set).
func (c *Client) GetAgent(ctx context.Context) (Agent, error) {
	id := c.AgentID
	q := url.Values{}
	q.Set("agentid", fmt.Sprintf("%d", id))
	var out agentGetResp
	if err := c.get(ctx, "/cgi-bin/agent/get", q, &out); err != nil {
		return Agent{}, err
	}
	users := make([]string, 0, len(out.AllowUserinfos.User))
	for _, u := range out.AllowUserinfos.User {
		users = append(users, u.UserID)
	}
	return Agent{
		AgentID: out.AgentID, Name: out.Name, SquareLogoURL: out.SquareLogoURL,
		Description: out.Description, Close: out.Close, RedirectDomain: out.RedirectDomain,
		ReportLocationFlag: out.ReportLocationFlag, IsReportEnter: out.IsReportEnter,
		HomeURL: out.HomeURL, CustomizedPublishStatus: out.CustomizedPublishStatus,
		AllowUserIDs: users, AllowPartyIDs: out.AllowPartys.PartyID, AllowTagIDs: out.AllowTags.TagID,
	}, nil
}

// AgentSummary is one row from agent/list.
type AgentSummary struct {
	AgentID       int    `json:"agentid"`
	Name          string `json:"name"`
	SquareLogoURL string `json:"square_logo_url"`
}

// ListAgents returns applications the current token can access.
func (c *Client) ListAgents(ctx context.Context) ([]AgentSummary, error) {
	var out struct {
		apiMeta
		AgentList []AgentSummary `json:"agentlist"`
	}
	if err := c.get(ctx, "/cgi-bin/agent/list", nil, &out); err != nil {
		return nil, err
	}
	return out.AgentList, nil
}

// AgentPatch is the body for agent/set. Zero values are omitted.
type AgentPatch struct {
	AgentID            int    `json:"agentid"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	LogoMediaID        string `json:"logo_mediaid,omitempty"`
	RedirectDomain     string `json:"redirect_domain,omitempty"`
	HomeURL            string `json:"home_url,omitempty"`
	ReportLocationFlag *int   `json:"report_location_flag,omitempty"`
	IsReportEnter      *int   `json:"isreportenter,omitempty"`
}

// SetAgent updates application settings.
func (c *Client) SetAgent(ctx context.Context, patch AgentPatch) error {
	if patch.AgentID == 0 {
		patch.AgentID = c.AgentID
	}
	return c.post(ctx, "/cgi-bin/agent/set", patch, nil)
}
