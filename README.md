# ST: sharetrace


## 管理界面
#### 登录账号
* URL :<http://st.apptao.com/web>
* 账号:guest@appwill.com
* 密码:123456

#### 效果图
 ![image](https://github.com/appwilldev/sharetrace/blob/master/web/img/stat_demo.png)
---

## 记录分享
#### 分享通知,告诉ST分享链接的信息
* URL:  /1/st/share
* POST: 

```
* share_url: 必填，分享链接，最好包含3个信息<appid, fromid, itemid>
* fromid:    必填，分享的用户ID, 可以用IDFA，或者注册后的UserID，但是不要混用
* appid:     必填，App的ID，需要在管理系统中注册并录入该AppID
* itemid:    选填，分享的内容ID
* channel:   选填，weixin,qq,weibo,else
* ver:       选填，App的版本
* des:       选填，备注信息
```

* Example: 
```
      curl -l -H "Content-type: application/json"       -X POST        -d '{"fromid":"8931","appid":"123", "itemid":"d1ff1099f8J", "share_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934","appname":"AV_FunnyTime", "channel":"weixin", "des":"test"}'            "http://localhost:8580/1/st/share"

```

---

## 跟踪点击 
#### 把下面代码加入分享页面尾部(不要加在head中, 放在body标签后面) 
* Code:
```

<script src="http://st.apptao.com/web/js/sharetrace.js"></script>

```

* 这段代码会把iFrame, 嵌入分享页面，便于在ST域名下跟踪Cookie，以及用户IP
* 如果是在Safari打开的页面，会在ST的域名(st.apptao.com)下面种下stcookieid，作为跟踪标识
* 如果是在微信／QQ打开的页面，种下的 stcookieid，无法跟踪，只能采取通过 IP 的方式跟踪


---

## 跟踪安装

#### 客户端打开ST的指定页面
* 这个动作不会对客户端体验有任何影响，不用跳出App到Safari，可以理解为在后台打开一个View，然后自动关闭了
* 原理是采用SFSafariViewController 打开一个透明的View，ST会从URL获取ST域名下面的stcookieid
 
* URL: http://st.apptao.com/1/st/webbeaconcheck?appid=???&installid=???
```
* appid:必填, 就是用户在管理系统中配置的appid，此参数必须有
* installid: 必填, 用户ID, 可以用IDFA，或者注册后的UserID，但是不要混用
```

* 客户端代码可以参考：<https://github.com/mackuba/SafariAutoLoginTest>

* 这个动作应该在App刚刚安装的时候做，以后就不用再执行这个检查了

