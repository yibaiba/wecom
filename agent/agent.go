package agent

import (
	"context"
	"net/url"
	"strconv"
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

// Get returns the current application.
func (a *API) Get(ctx context.Context) (Agent, error) {
	q := url.Values{}
	q.Set("agentid", strconv.Itoa(a.App.AgentID))
	var out agentGetResp
	if err := a.App.GetJSON(ctx, "/cgi-bin/agent/get", q, &out); err != nil {
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

// Summary is one row from agent/list.
type Summary struct {
	AgentID       int    `json:"agentid"`
	Name          string `json:"name"`
	SquareLogoURL string `json:"square_logo_url"`
}

// List returns applications the current token can access.
func (a *API) List(ctx context.Context) ([]Summary, error) {
	var out struct {
		AgentList []Summary `json:"agentlist"`
	}
	if err := a.App.GetJSON(ctx, "/cgi-bin/agent/list", nil, &out); err != nil {
		return nil, err
	}
	return out.AgentList, nil
}

// Patch is the body for agent/set. Zero values are omitted.
type Patch struct {
	AgentID            int    `json:"agentid"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	LogoMediaID        string `json:"logo_mediaid,omitempty"`
	RedirectDomain     string `json:"redirect_domain,omitempty"`
	HomeURL            string `json:"home_url,omitempty"`
	ReportLocationFlag *int   `json:"report_location_flag,omitempty"`
	IsReportEnter      *int   `json:"isreportenter,omitempty"`
}

// Set updates application settings.
func (a *API) Set(ctx context.Context, patch Patch) error {
	if patch.AgentID == 0 {
		patch.AgentID = a.App.AgentID
	}
	return a.App.PostJSON(ctx, "/cgi-bin/agent/set", patch, nil)
}
