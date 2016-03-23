package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
	"github.com/gin-gonic/gin"
)

func AppUserMoney(c *gin.Context) {
	q := c.Request.URL.Query()
	userIdStr := q["userid"][0]
	appIdStr := q["appid"][0]

	appDB, err := models.GetAppInfoByAppid(nil, appIdStr)
	if err != nil {
		Error(c, DATA_NOT_FOUND, nil, nil)
		return
	}

	award_str := ""
	if appDB.Status == 0 || appDB.Yue < 1000 {
		award_str = "分享暂时没有奖励规则哦, 请您继续关注!"
	} else {
		//if appDB.ShareClickMoney > 0 {
		//	award_str = fmt.Sprintf("%s分享获得点击, 每次奖励分享者%.2f元, 目前最多还可以有%d份奖励;", award_str, float64(appDB.ShareClickMoney/100), appDB.Yue/appDB.ShareClickMoney)
		//}
		if appDB.ShareInstallMoney > 0 && appDB.InstallMoney > 0 {
			award_str = fmt.Sprintf("%s分享获得安装, 奖励分享者%.2f元, 安装者%.2f元，还有%d份奖励", award_str, float64(appDB.ShareInstallMoney)/100.0, float64(appDB.InstallMoney)/100.0, appDB.Yue/(appDB.ShareInstallMoney+appDB.InstallMoney))
		} else {
			if appDB.ShareInstallMoney > 0 {
				award_str = fmt.Sprintf("%s分享获得安装, 奖励分享者%.2f元, 还有%d份奖励", award_str, float64(appDB.ShareInstallMoney)/100.0, appDB.Yue/appDB.ShareInstallMoney)
			}
			if appDB.InstallMoney > 0 {
				award_str = fmt.Sprintf("%s分享获得安装, 奖励安装者%.2f元, 还有%d份奖励", award_str, float64(appDB.InstallMoney)/100.0, appDB.Yue/appDB.InstallMoney)
			}
		}
	}

	dataList, _, _ := models.GetAppuserMoneyListByUserid(nil, appIdStr, userIdStr)
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["award_str"] = award_str
	var total = float64(0.0)

	time_now := time.Now()
	year, month, day := time_now.Date()
	date_now := fmt.Sprintf("%d-%d-%d", year, month, day)
	total_today := 0.0
	used := 0.0
	for _, row := range dataList {
		created_utc := time.Unix(int64(row.CreatedUTC), 0)
		year, mon, day := created_utc.Date()
		hour, min, sec := created_utc.Clock()
		date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)

		if row.MoneyType == conf.MONEY_TYPE_HFCZ {
			row.Money = -float64(row.Money / 100.0)
			used = used + row.Money
		} else {
			row.Money = float64(row.Money / 100.0)
			if date_tmp == date_now {
				total_today = total_today + float64(row.Money)
			}
			total = total + float64(row.Money)
		}

		s := fmt.Sprintf("%d-%d-%d %02d:%02d:%02d\n", year, mon, day, hour, min, sec)
		row.Des = row.Des + "      " + s
	}
	data["appid"] = appIdStr
	data["appuserid"] = userIdStr
	data["total_today"] = fmt.Sprintf("%.2f", total_today)
	data["total"] = fmt.Sprintf("%.2f", total)
	data["total_left"] = fmt.Sprintf("%.2f", total+used)
	data["res"] = dataList

	c.HTML(200, "appusermoney.html", data)
	return
}

func AppUserScore(c *gin.Context) {
	q := c.Request.URL.Query()
	userIdStr := q["userid"][0]
	appIdStr := q["appid"][0]

	appDB, err := models.GetAppInfoByAppid(nil, appIdStr)
	if err != nil {
		Error(c, DATA_NOT_FOUND, nil, nil)
		return
	}

	award_str := "100积分=1元;"
	if appDB.Status == 0 {
		award_str = ""
	} else {
		if appDB.ShareClickMoney > 0 {
			award_str = fmt.Sprintf("%s分享获得点击, 奖励分享者%d分;", award_str, appDB.ShareClickMoney)
		}
		if appDB.ShareInstallMoney > 0 {
			award_str = fmt.Sprintf("%s分享获得安装, 奖励分享者%d分;", award_str, appDB.ShareInstallMoney)
		}
		if appDB.InstallMoney > 0 {
			award_str = fmt.Sprintf("%s安装者%d分;", award_str, appDB.ShareInstallMoney)
		}
	}

	dataList, _, _ := models.GetAppuserMoneyListByUserid(nil, appIdStr, userIdStr)
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["award_str"] = award_str
	var total = 0

	time_now := time.Now()
	year, month, day := time_now.Date()
	date_now := fmt.Sprintf("%d-%d-%d", year, month, day)
	total_today := 0
	used := 0
	for _, row := range dataList {
		created_utc := time.Unix(int64(row.CreatedUTC), 0)
		year, mon, day := created_utc.Date()
		hour, min, sec := created_utc.Clock()
		date_tmp := fmt.Sprintf("%d-%d-%d", year, month, day)

		if row.MoneyType == conf.MONEY_TYPE_HFCZ {
			row.Money = -row.Money
			used = used + int(row.Money)
		} else {
			row.Money = row.Money
			if date_tmp == date_now {
				total_today = total_today + int(row.Money)
			}
			total = total + int(row.Money)
		}

		s := fmt.Sprintf("%d-%d-%d %02d:%02d:%02d\n", year, mon, day, hour, min, sec)
		row.Des = row.Des + "      " + s
	}
	data["appid"] = appIdStr
	data["appuserid"] = userIdStr
	data["total_today"] = fmt.Sprintf("%d", total_today)
	data["total"] = fmt.Sprintf("%d", total)
	data["total_left"] = fmt.Sprintf("%d", total+used)
	data["res"] = dataList

	clientIP := c.ClientIP()
	md5Ctx := md5.New()
	sign_info := fmt.Sprintf("%s_%s_%s", appIdStr, userIdStr, clientIP)
	md5Ctx.Write([]byte(sign_info))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))
	data["sign"] = sign

	c.HTML(200, "appuserscore.html", data)
	return
}
