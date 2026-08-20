package wecom

import (
	"encoding/json"
	"testing"
)

func TestIdentityFromUserGetMapsOfficialFields(t *testing.T) {
	var out userGetResp
	if err := json.Unmarshal([]byte(officialUserGetJSON), &out); err != nil {
		t.Fatal(err)
	}
	got := identityFromUserGet("zhangsan", out)
	if got.Email != "zhangsan@qq.com" || got.BizMail != "zhangsan@tencent.com" {
		t.Fatalf("mail %+v", got)
	}
	if got.Avatar != "http://wx.qlogo.cn/mmopen/dir/0" || got.ThumbAvatar != "http://wx.qlogo.cn/mmopen/dir/100" {
		t.Fatalf("avatar %+v", got)
	}
	if got.QRCode == "" || got.Telephone != "020-123456" || got.Address != "广州市海珠区新港中路" {
		t.Fatalf("contact %+v", got)
	}
	if got.Status != 1 || got.MainDepartment != 1 || got.Gender != "1" {
		t.Fatalf("status %+v", got)
	}
	if len(got.Department) != 2 || got.DirectLeader[0] != "lisi" {
		t.Fatalf("org %+v", got)
	}
	if len(got.ExtAttr) != 2 || got.ExternalProfile.WechatChannels.Status != 1 {
		t.Fatalf("attrs %+v %+v", got.ExtAttr, got.ExternalProfile)
	}
}

func TestIdentityFromUserGetAcceptsNumericGender(t *testing.T) {
	var out userGetResp
	if err := json.Unmarshal([]byte(`{"gender":2,"name":"李四"}`), &out); err != nil {
		t.Fatal(err)
	}
	got := identityFromUserGet("lisi", out)
	if got.Gender != "2" {
		t.Fatalf("gender %q", got.Gender)
	}
}

func TestOverlayPrivateDetailFillsSensitiveGaps(t *testing.T) {
	got := overlayPrivateDetail(Identity{
		UserID: "zhangsan", Name: "目录名", Email: "dir@example.com", BizMail: "dir@corp.com",
	}, userDetailResp{
		UserID: "zhangsan", Name: "授权名", Gender: "1", Avatar: "https://a", QRCode: "https://q",
		Mobile: "138", Email: "a@b.com", BizMail: "a@corp.com", Address: "广州",
	})
	if got.Name != "授权名" || got.Email != "a@b.com" || got.BizMail != "a@corp.com" {
		t.Fatalf("overlay %+v", got)
	}
	if got.Mobile != "138" || got.Address != "广州" || got.QRCode != "https://q" || got.Gender != "1" {
		t.Fatalf("sensitive %+v", got)
	}
}

func TestOverlayPrivateDetailDoesNotFoldBizMailIntoEmail(t *testing.T) {
	got := overlayPrivateDetail(Identity{UserID: "u"}, userDetailResp{BizMail: "a@corp.com"})
	if got.Email != "" || got.BizMail != "a@corp.com" {
		t.Fatalf("fold %+v", got)
	}
}

func TestIdentityTicketsOmittedFromJSON(t *testing.T) {
	raw, err := json.Marshal(Identity{UserID: "zhangsan", UserTicket: "t", UserDocTicket: "d"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["user_ticket"]; ok {
		t.Fatalf("user_ticket present %s", raw)
	}
	if _, ok := m["user_doc_ticket"]; ok {
		t.Fatalf("user_doc_ticket present %s", raw)
	}
	if m["userid"] != "zhangsan" {
		t.Fatalf("json %s", raw)
	}
}
