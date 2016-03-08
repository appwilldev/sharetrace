# ST: sharetrace

分享跟踪系统，统计App用户的分享到社交网络的URL，跟踪、分析社交网络上的点击，以及由此带来的App新增用户

---

 ![image](https://github.com/appwilldev/sharetrace/blob/master/web/img/stat_demo.png)

---

## 创建账号
#### 注册管理账号，录入App相关信息
* URL : <http://st.apptao.com/web>
```
测试账号: guest@appwill.com,  密码:123456
```
---

## 记录分享
#### 用户分享后，通知ST系统分享链接的信息
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
#### 跟踪记录URL被点击的情况
* 把下面代码加入分享页面尾部(不要加在head中, 放在body标签后面) 
* Code:
```
<script src="http://st.apptao.com/web/js/sharetrace.js"></script>
```
* 如果是在Safari打开的页面，会在ST的域名(st.apptao.com)下面种下stcookieid，作为跟踪标识
* 如果是在微信／QQ打开的页面，种下的 stcookieid，无法跟踪，只能采取通过 IP 的方式跟踪
* 如果要跟踪页面上按钮点击情况, 可以在跳转到AppStore前调用, buttonid用来区分不同buton的点击, 不传则默认为1：
```
gotoAppStore(buttonid);
```

---

## 跟踪安装
#### 客户端"隐式"打开ST的指定链接，对App用户无任何影响
* 可以理解为在App后台打开一个透明的View，然后自动关闭了，建议在App刚刚安装的时候做一次即可
* 原理是采用SFSafariViewController 打开一个透明的View，ST会从URL获取ST域名下面的stcookieid
* URL: http://st.apptao.com/1/st/webbeaconcheck?appid=???&installid=???
```
* appid:必填, 就是用户在管理系统中配置的appid，此参数必须有
* installid: 必填, 用户ID, 可以用IDFA，或者注册后的UserID，但是不要混用
```
* 客户端代码可以参考：<https://github.com/mackuba/SafariAutoLoginTest>
* Obj-C 核心代码参考： <https://github.com/appwilldev/sharetrace/blob/master/web/js/objcexam.m>

