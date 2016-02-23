curl -l -H "Content-type: application/json" -X POST -d '{"email":"leo@appwill.com","passwd":"123456", "name":"leo"}' "http://localhost:8580/1/user/register"

curl -l -H "Content-type: application/json" -X POST -d '{"email":"leo@appwill.com","passwd":"123456"}' -D cookie.txt "http://localhost:8580/1/user/login"

curl -l -H "Content-type: application/json" -X POST -b cookie.txt "http://localhost:8580/1/user/logout"

curl -l -H "Content-type: application/json" -X POST -b cookie.txt -d '{"from":"8931","code":"0410a4e72689d065448a7d1ff1099f8J","link_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934","app":"AV_FunnyTime","installed":0}'  "http://localhost:8580/1/user/weixinlink/click"

curl -l -H "Content-type: application/json" -X POST  -d '{"fromid":"8931","appid":"123", "itemid":"d1ff1099f8J","share_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934","appname":"AV_FunnyTime", "channel":"weixin", "des":"test"}' "http://localhost:8580/1/st/share"

curl -l -H "Content-type: application/json" -X POST  -d '{"share_url":"http://appforvideo.com/v/8923?app=AV_FunnyTime&from=8934"}' "http://localhost:8580/1/st/click"

curl -l -H "Content-type: application/json" -X POST  -d '{"st_cookieid":"st_2016_2016","installid":"aaaaaaa"}' "http://localhost:8580/1/st/install"

curl -l -H "Content-type: application/json" "http://localhost:8580/1/st/score?userid=8931&appid=123"
