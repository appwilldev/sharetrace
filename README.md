# ST: sharetrace

### 分享
    #### 分享通知,告诉ST分享链接的信息
    * URL:  /1/st/share
    * POST: 
          * share_url: 必填，分享链接，最好包含3个信息<appid, fromid, itemid>
          * fromid:    必填，IDFA 
          * appid:     必填，用户注册后，新增的AppID
          * itemid:     分享内容ID
          * channel:    weixin,qq,weibo,else
          * ver:        App版本
          * des:        备注
    * Example: curl -l -H "Content-type: application/json" -X POST  -d '{"fromid":"8931","appid":"123", "itemid":"d1ff1099f8J","share_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934","appname":"AV_FunnyTime", "channel":"weixin", "des":"test"}' "http://localhost:8580/1/st/share"


### 跟踪 
    #### 分享页面iFrame, 嵌入分享页面，便于ST种下跟踪Cookie，以及跟踪用户IP
    * Code: <iframe src="http://st.appforvideo.com/1/st/webbeacon?share_url=encoded_url"></iframe>
    * Example: <iframe src="http://192.168.1.17:8580/1/st/webbeacon?share_url=http%3A%2F%2Fappforvideo.com%2Fv%2F8923%3Fapp%3DAV_FunnyTime%26from%3D8934"></iframe>
    * 如果是在Safari打开的页面，会在ST的域名(st.appforvideo.com)下面种下stcookieid，作为跟踪标识
    * 如果是在微信／QQ打开的页面，种下的stcookieid，无法跟踪，只能采取通过IP的方式跟踪



### 安装
    #### 客户端打开ST的指定页面，获取stcookieid，作为trackid
      * 采用SFSafariViewController 可以在不影响用户操作的情况下，偷偷的打开一个View，从下面的URL获取ST域名下面的stcookieid
      * URL: http://st.appforvideo.com/1/st/webbeaconcheck?appid=??? (appid就是用户在管理系统中配置的appid，此参数必须有)
      * 客户端代码可以参考：https://github.com/mackuba/SafariAutoLoginTest
    #### 如果没有stcookieid则使用IP作为trackid，且click_type = 0
    #### 安装通知,告诉ST新安装App设备的信息
      * URL:  /1/st/install
      * POST: 
            * click_type:必填 
                * 0:使用cookieid，
                * 1:用IP
            * installid: 必填,IDFA(451141EE-9540-4AA1-A82A-350AC548193D) 
            * trackid: 必填,为SFSafariViewController获得的stcookieid； 当click_type=1时为客户端的真实IP(比如：58.83.170.249)
      * Example: curl -l -H "Content-type: application/json" -X POST  -d '{"trackid":"st_2016_2016","click_type":"0","installid":"451141EE-9540-4AA1-A82A-350AC548193D"}' "http://localhost:8580/1/st/install"

