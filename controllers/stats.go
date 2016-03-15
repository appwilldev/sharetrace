package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/appwilldev/sharetrace/models"
	"github.com/gin-gonic/gin"
)

func StatsShare(c *gin.Context) {
	log.Println("")
	q := c.Request.URL.Query()

	appIdStr := q["appid"][0]

	dataList, _ := models.GetShareClickListByAppid(nil, appIdStr)
	var data map[string]interface{}
	data = make(map[string]interface{})
	var total = float64(0.0)
	for _, row := range dataList {
		row.ClickSession.Des = "用户：" + row.ClickSession.Cookieid + " 点击了用户ID:" + row.ShareURL.Fromid + "的分享链接:" + row.ShareURL.ShareURL
		if row.ClickSession.Installid != "" {
			row.ClickSession.Des = "推荐下载 获得100分: " + row.ClickSession.Des
			total = total + 100
		} else {
			row.ClickSession.Des = "分享被点击 获得1分: " + row.ClickSession.Des
			total = total + 1
		}
		created_utc := time.Unix(int64(row.ClickSession.CreatedUTC), 0)
		year, mon, day := created_utc.Date()
		hour, min, sec := created_utc.Clock()
		s := fmt.Sprintf("%d-%d-%d %02d:%02d:%02d\n", year, mon, day, hour, min, sec)
		row.ClickSession.Des = row.ClickSession.Des + s
		row.ScoreDes = row.ClickSession.Des

	}
	data["total"] = total
	data["res"] = dataList

	ret := gin.H{"status": true}
	ret["data"] = data
	c.JSON(200, ret)
	return
}

func StatsTotal(c *gin.Context) {

	q := c.Request.URL.Query()
	appIdStr := q["appid"][0]

	var delta int
	delta = 7
	time_now := time.Now().AddDate(0, 0, -delta+1)

	var t_start int64
	if len(q["start"]) > 0 {
		t_start_str := q["start"][0]
		t_start, _ = strconv.ParseInt(t_start_str, 10, 64)
		start_utc := time.Unix(t_start, 0)
		time_now = start_utc
		//year, month, day := start_utc.Date()
		//date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)
	}

	var t_end int64
	if len(q["end"]) > 0 {
		t_end_str := q["end"][0]
		t_end, _ = strconv.ParseInt(t_end_str, 10, 64)
		//end_utc := time.Unix(t_end, 0)
		//year, month, day := end_utc.Date()
		//date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)
	}

	if t_start > 0 && t_end > 0 {
		delta = int((t_end - t_start) / (24 * 3600))
	}

	ret := gin.H{"status": true}

	var data map[string]interface{}
	data = make(map[string]interface{})

	share := make([]interface{}, 0)
	click := make([]interface{}, 0)
	button := make([]interface{}, 0)
	install := make([]interface{}, 0)

	for i := 0; i < delta; i++ {
		time_now_tmp := time_now.AddDate(0, 0, +i)
		year, month, day := time_now_tmp.Date()
		date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)

		share_tmp := make(map[string]interface{}, 1)
		share_total, _ := models.GetShareTotalByAppid(nil, appIdStr, date_tmp)
		share_tmp[date_tmp] = share_total
		share = append(share, share_tmp)

		click_tmp := make(map[string]interface{}, 1)
		click_total, _ := models.GetClickTotalByAppid(nil, appIdStr, date_tmp)
		click_tmp[date_tmp] = click_total
		click = append(click, click_tmp)

		button_tmp := make(map[string]interface{}, 1)
		button_total, _ := models.GetButtonTotalByAppid(nil, appIdStr, date_tmp)
		button_tmp[date_tmp] = button_total
		button = append(button, button_tmp)

		install_tmp := make(map[string]interface{}, 1)
		install_total, _ := models.GetInstallTotalByAppid(nil, appIdStr, date_tmp)
		install_tmp[date_tmp] = install_total
		install = append(install, install_tmp)
	}

	data["share"] = share
	data["click"] = click
	data["install"] = install
	data["button"] = button

	ret["data"] = data
	c.JSON(200, ret)
}

func StatsHost(c *gin.Context) {
	q := c.Request.URL.Query()
	host := q["host"][0]

	var delta int
	delta = 7
	time_now := time.Now().AddDate(0, 0, -delta+1)

	var t_start int64
	if len(q["start"]) > 0 {
		t_start_str := q["start"][0]
		t_start, _ = strconv.ParseInt(t_start_str, 10, 64)
		start_utc := time.Unix(t_start, 0)
		time_now = start_utc
	}

	var t_end int64
	if len(q["end"]) > 0 {
		t_end_str := q["end"][0]
		t_end, _ = strconv.ParseInt(t_end_str, 10, 64)
	}

	if t_start > 0 && t_end > 0 {
		delta = int((t_end - t_start) / (24 * 3600))
	}

	ret := gin.H{"status": true}

	var data map[string]interface{}
	data = make(map[string]interface{})

	click := make([]interface{}, 0)
	button := make([]interface{}, 0)

	iphone := make([]interface{}, 0)
	ipad := make([]interface{}, 0)
	android := make([]interface{}, 0)
	window := make([]interface{}, 0)
	phone_else := make([]interface{}, 0)

	safari := make([]interface{}, 0)
	wechat := make([]interface{}, 0)
	qq := make([]interface{}, 0)
	weibo := make([]interface{}, 0)
	chrome := make([]interface{}, 0)
	browser_else := make([]interface{}, 0)

	for i := 0; i < delta; i++ {
		time_now_tmp := time_now.AddDate(0, 0, +i)
		year, month, day := time_now_tmp.Date()
		date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)

		click_tmp := make(map[string]interface{}, 1)
		click_total, _ := models.GetClickTotalByHost(nil, host, date_tmp)
		click_tmp[date_tmp] = click_total
		click = append(click, click_tmp)

		button_tmp := make(map[string]interface{}, 1)
		button_total, _ := models.GetButtonTotalByHost(nil, host, date_tmp)
		button_tmp[date_tmp] = button_total
		button = append(button, button_tmp)

		iphone_tmp := make(map[string]interface{}, 1)
		iphone_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "iPhone OS")
		iphone_tmp[date_tmp] = iphone_total
		iphone = append(iphone, iphone_tmp)

		ipad_tmp := make(map[string]interface{}, 1)
		ipad_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "iPad")
		ipad_tmp[date_tmp] = ipad_total
		ipad = append(ipad, ipad_tmp)

		android_tmp := make(map[string]interface{}, 1)
		android_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "Android")
		android_tmp[date_tmp] = android_total
		android = append(android, android_tmp)

		window_tmp := make(map[string]interface{}, 1)
		window_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "Window")
		window_tmp[date_tmp] = window_total
		window = append(window, window_tmp)

		phone_else_tmp := make(map[string]interface{}, 1)
		phone_else_total := click_total - iphone_total - ipad_total - android_total - window_total
		phone_else_tmp[date_tmp] = phone_else_total
		phone_else = append(phone_else, phone_else_tmp)

		safari_tmp := make(map[string]interface{}, 1)
		safari_total, _ := models.GetTotalByHostiPhone(nil, host, date_tmp, "Safari")
		safari_tmp[date_tmp] = safari_total
		safari = append(safari, safari_tmp)

		wechat_tmp := make(map[string]interface{}, 1)
		wechat_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "MicroMessenger")
		wechat_tmp[date_tmp] = wechat_total
		wechat = append(wechat, wechat_tmp)

		qq_tmp := make(map[string]interface{}, 1)
		qq_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "QQ")
		qq_tmp[date_tmp] = qq_total
		qq = append(qq, qq_tmp)

		weibo_tmp := make(map[string]interface{}, 1)
		weibo_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "Weibo")
		weibo_tmp[date_tmp] = weibo_total
		weibo = append(weibo, weibo_tmp)

		chrome_tmp := make(map[string]interface{}, 1)
		chrome_total, _ := models.GetTotalByHostPhone(nil, host, date_tmp, "Chrome")
		chrome_tmp[date_tmp] = chrome_total
		chrome = append(chrome, chrome_tmp)

		browser_else_tmp := make(map[string]interface{}, 1)
		browser_else_total := click_total - safari_total - wechat_total - qq_total - weibo_total
		browser_else_tmp[date_tmp] = browser_else_total
		browser_else = append(browser_else, browser_else_tmp)

	}
	data["click"] = click
	data["button"] = button

	data["iphone"] = iphone
	data["ipad"] = ipad
	data["android"] = android
	data["window"] = window
	data["phone_else"] = phone_else

	data["safari"] = safari
	data["wechat"] = wechat
	data["qq"] = qq
	data["weibo"] = weibo
	data["chrome"] = chrome
	data["browser_else"] = browser_else

	ret["data"] = data
	c.JSON(200, ret)
}
