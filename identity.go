package wecom

// Identity is the member or visitor after exchanging a WeCom authorization code.
// Values come from user/get, auth/getuserinfo, and auth/getuserdetail when
// those APIs return them. Empty means the corp app or the member did not grant
// that field. Email is the personal mailbox; BizMail is the corp mailbox.
type Identity struct {
	UserID           string `json:"userid,omitempty"`
	Name             string `json:"name,omitempty"`
	Alias            string `json:"alias,omitempty"`
	Position         string `json:"position,omitempty"`
	ExternalPosition string `json:"external_position,omitempty"`
	Email            string `json:"email,omitempty"`
	BizMail          string `json:"biz_mail,omitempty"`
	Mobile           string `json:"mobile,omitempty"`
	Telephone        string `json:"telephone,omitempty"`
	Address          string `json:"address,omitempty"`
	// Gender is "0" undefined, "1" male, "2" female.
	Gender      string `json:"gender,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	ThumbAvatar string `json:"thumb_avatar,omitempty"`
	QRCode      string `json:"qr_code,omitempty"`
	// Status is 1 active, 2 disabled, 4 inactive, 5 quit.
	Status          int             `json:"status,omitempty"`
	MainDepartment  int             `json:"main_department,omitempty"`
	Department      []int           `json:"department,omitempty"`
	Order           []int           `json:"order,omitempty"`
	IsLeaderInDept  []int           `json:"is_leader_in_dept,omitempty"`
	DirectLeader    []string        `json:"direct_leader,omitempty"`
	OpenUserID      string          `json:"open_userid,omitempty"`
	DeviceID        string          `json:"device_id,omitempty"`
	UserTicket      string          `json:"-"`
	UserDocTicket   string          `json:"-"`
	OpenID          string          `json:"openid,omitempty"`
	ExternalUserID  string          `json:"external_userid,omitempty"`
	ExtAttr         []ExtAttr       `json:"extattr,omitempty"`
	ExternalProfile ExternalProfile `json:"external_profile,omitempty"`
}

// ExtAttr is a directory or external custom attribute.
type ExtAttr struct {
	Type        int                `json:"type"`
	Name        string             `json:"name"`
	Text        ExtAttrText        `json:"text"`
	Web         ExtAttrWeb         `json:"web"`
	Miniprogram ExtAttrMiniprogram `json:"miniprogram"`
}

// ExtAttrText is type 0.
type ExtAttrText struct {
	Value string `json:"value"`
}

// ExtAttrWeb is type 1.
type ExtAttrWeb struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// ExtAttrMiniprogram is type 2.
type ExtAttrMiniprogram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
	Title    string `json:"title"`
}

// ExtAttrGroup is the extattr object on user/create.
type ExtAttrGroup struct {
	Attrs []ExtAttr `json:"attrs"`
}

// ExternalProfile is the member's customer-facing profile.
type ExternalProfile struct {
	ExternalCorpName string         `json:"external_corp_name,omitempty"`
	WechatChannels   WechatChannels `json:"wechat_channels,omitempty"`
	ExternalAttr     []ExtAttr      `json:"external_attr,omitempty"`
}

// WechatChannels is the video-account card on ExternalProfile.
type WechatChannels struct {
	Nickname string `json:"nickname,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// PreferredEmail is the corp mailbox, else the personal mailbox. Empty if neither was granted.
func PreferredEmail(id Identity) string {
	return firstNonEmpty(id.BizMail, id.Email)
}

// MergeDirectory copies non-empty directory fields onto a login identity.
func MergeDirectory(login, dir Identity) Identity {
	out := login
	if uid := firstNonEmpty(dir.UserID); uid != "" {
		out.UserID = uid
	}
	out.Name = firstNonEmpty(dir.Name, login.Name)
	out.Email = firstNonEmpty(dir.Email, login.Email)
	out.BizMail = firstNonEmpty(dir.BizMail, login.BizMail)
	out.Avatar = firstNonEmpty(dir.Avatar, login.Avatar)
	out.ThumbAvatar = firstNonEmpty(dir.ThumbAvatar, login.ThumbAvatar)
	return out
}

// PreferredAvatar is the full avatar, else the thumbnail. Empty if neither was granted.
func PreferredAvatar(id Identity) string {
	return firstNonEmpty(id.Avatar, id.ThumbAvatar)
}
