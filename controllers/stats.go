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
		row.ScoreDesc = row.ClickSession.Des

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

	}
	data["click"] = click
	data["button"] = button

	ret["data"] = data
	c.JSON(200, ret)

}
