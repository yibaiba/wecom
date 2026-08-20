package wecom

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// WriteLoginPanel renders the official WWLogin component. callback is the
// host app's code receiver and must be same-origin when redirect_type=callback.
func WriteLoginPanel(w http.ResponseWriter, app LoginApp, state, callback string) {
	payload, err := json.Marshal(map[string]string{
		"appid":        app.CorpID,
		"agentid":      strconv.Itoa(app.AgentID),
		"redirect_uri": callback,
		"state":        state,
		"callback":     callback,
	})
	if err != nil || !app.Configured() {
		http.Error(w, "wecom login is not configured", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	page := strings.Replace(loginPanelHTML, "__WECOM_JSON__", string(payload), 1)
	_, _ = w.Write([]byte(page))
}

const loginPanelHTML = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>企业微信登录</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:0;background:#f5f5f5;color:#111}
.wrap{max-width:480px;margin:48px auto;padding:24px;background:#fff;border-radius:12px;box-shadow:0 1px 4px rgba(0,0,0,.06)}
h1{font-size:20px;margin:0 0 8px}
p{color:#666;font-size:14px;line-height:1.5}
#panel{min-height:416px;display:flex;justify-content:center}
#panel iframe{display:block;border:0}
.err{color:#c00}
</style>
</head>
<body>
<div class="wrap">
<h1>企业微信登录</h1>
<p id="status">正在创建企业微信登录组件。</p>
<div id="panel"></div>
</div>
<script>window.__WECOM__=__WECOM_JSON__;</script>
<script>
(function(){
  var cfg=window.__WECOM__||{};
  var status=document.getElementById("status");
  var panel=document.getElementById("panel");
  function go(code){
    var u=new URL(cfg.callback, window.location.origin);
    u.searchParams.set("code", code);
    u.searchParams.set("state", cfg.state);
    window.location.assign(u.toString());
  }
  function fail(msg){
    status.className="err";
    status.textContent=msg;
  }
  function mount(ww){
    ww.createWWLoginPanel({
      el:panel,
      params:{
        login_type:"CorpApp",
        appid:cfg.appid,
        agentid:cfg.agentid,
        redirect_uri:cfg.redirect_uri,
        state:cfg.state,
        redirect_type:"callback",
        panel_size:"middle",
        lang:"zh"
      },
      onCheckWeComLogin:function(ev){
        status.className="";
        status.textContent=ev.isWeComLogin
          ?"已检测到企业微信客户端，请在官方登录组件中确认。"
          :"未检测到企业微信客户端，请在官方登录组件中扫码登录。";
      },
      onOpenInWecom:function(){
        status.className="";
        status.textContent="已检测到企业微信客户端，请在官方登录组件中确认。";
      },
      onLoginSuccess:function(res){ go(res.code); },
      onLoginFail:function(err){ fail((err&&err.errMsg)||"企业微信登录组件授权失败"); }
    });
  }
  var s=document.createElement("script");
  s.src="https://wwcdn.weixin.qq.com/node/open/js/wecom-jssdk-2.3.4.js";
  s.async=true;
  s.onload=function(){
    if(!window.ww||!window.ww.createWWLoginPanel){
      fail("企业微信登录组件加载失败");
      return;
    }
    mount(window.ww);
  };
  s.onerror=function(){ fail("企业微信登录组件加载失败"); };
  document.head.appendChild(s);
})();
</script>
</body>
</html>
`
