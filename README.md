# ST: sharetrace

### 分享通知,告诉ST分享链接的信息
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


### 分享页面iFrame, 嵌入分享页面，便于ST种下跟踪Cookie，以及跟踪用户IP
    * Code: <iframe src="http://st.appforvideo.com/1/st/webbeacon?share_url=encoded_url"></iframe>
    * Example: <iframe src="http://192.168.1.17:8580/1/st/webbeacon?share_url=http%3A%2F%2Fappforvideo.com%2Fv%2F8923%3Fapp%3DAV_FunnyTime%26from%3D8934"></iframe>



### 安装通知,告诉ST新安装App设备的信息

  #### 有stcookieid：
      * URL:  /1/st/install
      * POST: 
            * stcookieid: string
            * installid: string, IDFA 
      * Example: curl -l -H "Content-type: application/json" -X POST  -d '{"stcookieid":"st_2016_2016","installid":"aaaaaaa"}' "http://localhost:8580/1/st/install"
   #### 无stcookieid：
