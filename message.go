package wecom

import (
	"context"
	"encoding/json"
)

// Message is a cgi-bin/message/send payload. Fill MsgType plus the matching body.
type Message struct {
	ToUser                 string          `json:"touser,omitempty"`
	ToParty                string          `json:"toparty,omitempty"`
	ToTag                  string          `json:"totag,omitempty"`
	MsgType                string          `json:"msgtype"`
	AgentID                int             `json:"agentid"`
	Text                   *TextBody       `json:"text,omitempty"`
	Image                  *MediaBody      `json:"image,omitempty"`
	Voice                  *MediaBody      `json:"voice,omitempty"`
	Video                  *VideoBody      `json:"video,omitempty"`
	File                   *MediaBody      `json:"file,omitempty"`
	TextCard               *TextCardBody   `json:"textcard,omitempty"`
	News                   *NewsBody       `json:"news,omitempty"`
	MPNews                 *MPNewsBody     `json:"mpnews,omitempty"`
	Markdown               *TextBody       `json:"markdown,omitempty"`
	MiniprogramNotice      *MiniNoticeBody `json:"miniprogram_notice,omitempty"`
	TemplateCard           json.RawMessage `json:"template_card,omitempty"`
	Safe                   int             `json:"safe,omitempty"`
	EnableIDTrans          int             `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int             `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int             `json:"duplicate_check_interval,omitempty"`
}

// TextBody is text or markdown content.
type TextBody struct {
	Content string `json:"content"`
}

// MediaBody is image/voice/file.
type MediaBody struct {
	MediaID string `json:"media_id"`
}

// VideoBody is a video message.
type VideoBody struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// TextCardBody is a text card.
type TextCardBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

// NewsBody is an articles list.
type NewsBody struct {
	Articles []NewsArticle `json:"articles"`
}

// NewsArticle is one news item.
type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	PicURL      string `json:"picurl,omitempty"`
	AppID       string `json:"appid,omitempty"`
	PagePath    string `json:"pagepath,omitempty"`
}

// MPNewsBody is mpnews stored inside WeCom.
type MPNewsBody struct {
	Articles []MPNewsArticle `json:"articles"`
}

// MPNewsArticle is one mpnews item.
type MPNewsArticle struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	Content          string `json:"content"`
	Digest           string `json:"digest,omitempty"`
}

// MiniNoticeBody is a miniprogram notice.
type MiniNoticeBody struct {
	AppID             string         `json:"appid"`
	Page              string         `json:"page,omitempty"`
	Title             string         `json:"title"`
	Description       string         `json:"description,omitempty"`
	EmphasisFirstItem bool           `json:"emphasis_first_item,omitempty"`
	ContentItem       []MiniNoticeKV `json:"content_item,omitempty"`
}

// MiniNoticeKV is one miniprogram notice row.
type MiniNoticeKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SendResult is the message/send response.
type SendResult struct {
	InvalidUser    string `json:"invaliduser"`
	InvalidParty   string `json:"invalidparty"`
	InvalidTag     string `json:"invalidtag"`
	UnlicensedUser string `json:"unlicenseduser"`
	MsgID          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}

// Send pushes an application message.
func (c *Client) Send(ctx context.Context, msg Message) (SendResult, error) {
	if msg.AgentID == 0 {
		msg.AgentID = c.AgentID
	}
	var out struct {
		apiMeta
		SendResult
	}
	if err := c.post(ctx, "/cgi-bin/message/send", msg, &out); err != nil {
		return SendResult{}, err
	}
	return out.SendResult, nil
}

// SendText sends a text message to users (pipe-separated userids, or "@all").
func (c *Client) SendText(ctx context.Context, toUser, content string) (SendResult, error) {
	return c.Send(ctx, Message{ToUser: toUser, MsgType: "text", Text: &TextBody{Content: content}})
}

// SendMarkdown sends a markdown message.
func (c *Client) SendMarkdown(ctx context.Context, toUser, content string) (SendResult, error) {
	return c.Send(ctx, Message{ToUser: toUser, MsgType: "markdown", Markdown: &TextBody{Content: content}})
}

// RecallMessage recalls a sent message by msgid.
func (c *Client) RecallMessage(ctx context.Context, msgid string) error {
	return c.post(ctx, "/cgi-bin/message/recall", map[string]string{"msgid": msgid}, nil)
}

// UpdateTemplateCard updates a template card using response_code from Send.
func (c *Client) UpdateTemplateCard(ctx context.Context, payload any) error {
	return c.post(ctx, "/cgi-bin/message/update_template_card", payload, nil)
}
