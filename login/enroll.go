package login

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// WritePhoneQRPage renders a phone-scan page. statusPath and continuePath are
// host routes that the desktop browser polls after the phone finishes OAuth.
func WritePhoneQRPage(w http.ResponseWriter, authorize, statusPath, continuePath string) {
	if statusPath == "" {
		statusPath = EnrollStatusPath
	}
	if continuePath == "" {
		continuePath = EnrollContinuePath
	}
	img, err := qrPNGDataURI(authorize)
	if err != nil {
		http.Error(w, "failed to render QR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	page := phoneQRHTML
	page = strings.Replace(page, "__WECOM_AUTH__", jsonString(authorize), 1)
	page = strings.Replace(page, "__QR_DATA_URI__", img, 1)
	page = strings.Replace(page, "__STATUS_PATH__", jsonString(statusPath), 1)
	page = strings.Replace(page, "__CONTINUE_PATH__", jsonString(continuePath), 1)
	_, _ = w.Write([]byte(page))
}

// WriteEnrollDonePage tells the phone browser to return to the desktop.
func WriteEnrollDonePage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(enrollDoneHTML))
}

func qrPNGDataURI(content string) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 260)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

func jsonString(v string) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `""`
	}
	return string(b)
}

const enrollDoneHTML = `<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><title>企业微信授权</title></head><body style="font-family:sans-serif;padding:32px;"><h2>企业微信授权</h2><p>授权已完成，请回到电脑页面继续。</p></body></html>`

const phoneQRHTML = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>企业微信登录</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:0;background:#f5f5f5;color:#111}
.wrap{max-width:480px;margin:48px auto;padding:24px;background:#fff;border-radius:12px;box-shadow:0 1px 4px rgba(0,0,0,.06);text-align:center}
h1{font-size:20px;margin:0 0 8px}
p{color:#666;font-size:14px;line-height:1.5}
img.qr{margin:16px auto;display:block;width:260px;height:260px}
</style>
</head>
<body>
<div class="wrap">
<h1>请用手机企业微信扫码</h1>
<p>系统里还没有这个人，需要手机扫码完成登录。</p>
<img class="qr" width="260" height="260" alt="企业微信授权二维码" data-authorize=__WECOM_AUTH__ src="__QR_DATA_URI__">
<p id="status">等待手机授权…</p>
</div>
<script>
(function(){
  function poll(){
    fetch(__STATUS_PATH__,{credentials:"same-origin"}).then(function(r){return r.json()}).then(function(j){
      if(j.status==="completed"){ location.assign(__CONTINUE_PATH__); return; }
      setTimeout(poll,2000);
    }).catch(function(){ setTimeout(poll,2000); });
  }
  poll();
})();
</script>
</body>
</html>
`
