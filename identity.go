package wecom

// Identity is the member or visitor after exchanging a WeCom authorization code.
// Values come from user/get, user/getuserinfo (and auth/getuserinfo), and
// auth/getuserdetail when those APIs return them. Empty means the corp app or
// the member did not grant that field.
type Identity struct {
	UserID           string
	Name             string
	Alias            string
	Position         string
	ExternalPosition string
	Email            string
	BizMail          string
	Mobile           string
	Telephone        string
	Address          string
	// Gender is "0" undefined, "1" male, "2" female.
	Gender      string
	Avatar      string
	ThumbAvatar string
	QRCode      string
	// Status is 1 active, 2 disabled, 4 inactive, 5 quit.
	Status          int
	MainDepartment  int
	Department      []int
	Order           []int
	IsLeaderInDept  []int
	DirectLeader    []string
	OpenUserID      string
	DeviceID        string
	UserTicket      string
	UserDocTicket   string
	OpenID          string
	ExternalUserID  string
	ExtAttr         []ExtAttr
	ExternalProfile ExternalProfile
}

// ExtAttr is a directory or external custom attribute.
type ExtAttr struct {
	Type        int
	Name        string
	Text        ExtAttrText
	Web         ExtAttrWeb
	Miniprogram ExtAttrMiniprogram
}

// ExtAttrText is type 0.
type ExtAttrText struct {
	Value string
}

// ExtAttrWeb is type 1.
type ExtAttrWeb struct {
	URL   string
	Title string
}

// ExtAttrMiniprogram is type 2.
type ExtAttrMiniprogram struct {
	AppID    string
	PagePath string
	Title    string
}

// ExternalProfile is the member's customer-facing profile.
type ExternalProfile struct {
	ExternalCorpName string
	WechatChannels   WechatChannels
	ExternalAttr     []ExtAttr
}

// WechatChannels is the video-account card on ExternalProfile.
type WechatChannels struct {
	Nickname string
	Status   int
}
