// user
curl -l -H "Content-type: application/json" -X POST -d '{"email":"leo@appwill.com","passwd":"123456", "name":"leo"}' "http://localhost:8580/op/user/init"
curl -l -H "Content-type: application/json" -X POST -d '{"email":"leo@appwill.com","passwd":"123456"}' -D cookie.txt "http://localhost:8580/op/user/login"
curl -l -H "Content-type: application/json" -X POST -b cookie.txt "http://localhost:8580/op/user/logout"
curl -l -H "Content-type: application/json" "http://localhost:8580/op/users/1/20"

// app
curl -l -H "Content-type: application/json" -X POST -b cookie.txt -d '{"appid":"123456", "appname":"apptest", "appicon":"iconurl"}' "http://localhost:8580/op/app/new"
curl -l -H "Content-type: application/json" "http://localhost:8580/op/apps/all/1/20"

// st
curl -l -H "Content-type: application/json" -X POST  -d '{"fromid":"8888","appid":"1042901066", "itemid":"13104","share_url":"http://appforvideo.com/v/13104?app=AV_FunnyTime&fromid=8888","appname":"AV_FunnyTime", "channel":"weixin", "des":"test"}' "http://localhost:8580/1/st/share"

curl -l -H "Content-type: application/json" "http://localhost:8580/1/st/score?userid=8931&appid=123"

curl "http://localhost:8580/1/st/webbeacon?share_url=http%3A%2F%2Fappforvideo.com%2Fv%2F8923%3Fapp%3DAV_FunnyTime%26from%3D8934"

curl "http://localhost:8580/1/st/webbeaconbutton?buttonid=12&share_url=http%3A%2F%2Fappforvideo.com%2Fv%2F8923%3Fapp%3DAV_FunnyTime%26from%3D8934"

curl "http://localhost:8580/1/st/webbeaconcheck?appid=1042901066&installid=ccccccc"


// stats
curl -l -H "Content-type: application/json" "http://localhost:8580/1/stats/share?appid=123"

// sql
SELECT count(*) FROM "share_url" INNER JOIN "click_session" ON share_url.id=click_session.shareid WHERE ((appid='123') AND (installid is not null)) AND (date(to_timestamp(click_session.created_utc))='2016-02-23') 


SELECT count(*), (date(to_timestamp(click_session.created_utc)))  FROM "share_url" INNER JOIN "click_session" ON share_url.id=click_session.shareid WHERE ((appid='123') AND (installid is not null)) group by date(to_timestamp(click_session.created_utc)) ;


SELECT count(*) FROM "click_session" WHERE (url_host='i.apptao.com') AND (date(to_timestamp(created_utc))='2016-3-10')


