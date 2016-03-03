# ST: sharetrace

---

## 分享
#### 分享通知,告诉ST分享链接的信息
* URL:  /1/st/share
* POST: 

```
* share_url: 必填，分享链接，最好包含3个信息<appid, fromid, itemid>
* fromid:    必填，分享的用户ID
* appid:     必填，App的ID，需要在管理系统中注册并录入该AppID
* itemid:    选填，分享的内容ID
* channel:   选填，weixin,qq,weibo,else
* ver:       选填，App的版本
* des:       选填，备注信息
```

* Example: 

      curl -l -H "Content-type: application/json"       -X POST        -d '{"fromid":"8931","appid":"123", "itemid":"d1ff1099f8J", "share_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934","appname":"AV_FunnyTime", "channel":"weixin", "des":"test"}'            "http://localhost:8580/1/st/share"


---

## 跟踪 
#### 把下面代码加入分享页面， 这段代码会把iFrame, 嵌入分享页面，便于在ST域名下跟踪Cookie，以及用户IP
* Code:
```

<script language="JavaScript">
function set_stiframe(){
    iframe = document.createElement('iframe');
    var stiframe_url="http://st.apptao.com/1/st/webbeacon?share_url=" + encodeURIComponent(window.location);
    iframe.src = stiframe_url;
    document.body.appendChild(iframe);
}

set_stiframe()
</script>

```
        
        
* 如果是在Safari打开的页面，会在ST的域名(st.appforvideo.com)下面种下stcookieid，作为跟踪标识
* 如果是在微信／QQ打开的页面，种下的 stcookieid，无法跟踪，只能采取通过 IP 的方式跟踪


---

## 安装

#### 客户端打开ST的指定页面，获取stcookieid，作为trackid
* 采用SFSafariViewController 可以打开一个透明的View，从下面的URL获取ST域名下面的stcookieid,然后关掉这个View
 
* URL: http://st.appforvideo.com/1/st/webbeaconcheck?appid=??? (appid就是用户在管理系统中配置的appid，此参数必须有)

* SFSafariViewController打开的URL会跟据App Schema（比如：avft）跳回来，并且带上stcookieid, 比如：avft://stcookieid/st_2016_2016, 其中st_2016_2016即stcookieid


* 客户端代码可以参考：<https://github.com/mackuba/SafariAutoLoginTest>


#### 安装通知,告诉ST新安装App设备的信息
* URL:  /1/st/install
* POST: 
```
* click_type:必填, 0:使用Cookie方式; 1:用IP方式
* installid: 必填,比如：IDFA(451141EE-9540-4AA1-A82A-350AC548193D) 
* trackid:   必填,即stcookieid; 如果没有stcookieid,则用客户端的真实IP(比如：58.83.170.249)
```
    
* Example: 

      curl -l -H "Content-type: application/json" -X POST  -d '{"trackid":"st_2016_2016","click_type":"0","installid":"451141EE-9540-4AA1-A82A-350AC548193D"}' "http://localhost:8580/1/st/install"

