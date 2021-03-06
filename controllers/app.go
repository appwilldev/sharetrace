package controllers

import (
	"log"

	"github.com/appwilldev/sharetrace/models"
	"github.com/appwilldev/sharetrace/models/caches"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/gin-gonic/gin"
)

type NewAppPostData struct {
	Appid             string  `json:"appid" binding:"required"`
	AppName           string  `json:"appname" binding:"required"`
	AppSchema         string  `json:"appschema" binding:"required"`
	AppHost           string  `json:"apphost" binding:"required"`
	Status            int     `json:"status "`
	ShareClickMoney   float64 `json:"share_click_money"`
	ShareInstallMoney float64 `json:"share_install_money"`
	InstallMoney      float64 `json:"install_money"`
	AppIcon           string  `json:"appicon"`
}

func NewApp(c *gin.Context) {
	var reqData NewAppPostData
	err := c.BindJSON(&reqData)
	if err != nil {
		Error(c, BAD_POST_DATA)
		return
	}

	userid := getUserIdFromContext(c)
	if userid <= 0 {
		Error(c, LOGIN_NEEDED, nil, nil)
	}

	old_appid, err := caches.GetAppInfoIdByAppid(reqData.Appid)
	if old_appid != "" {
		log.Println("register duplicated appid:", reqData.Appid)
		Error(c, DATA_DUPLICATED, nil, nil)
		return
	}

	appDB, err := models.GetAppInfoByAppid(nil, reqData.Appid)
	if appDB != nil {
		log.Println("register duplicated appid:", reqData.Appid)
		Error(c, DATA_DUPLICATED, nil, nil)
		return
	}

	appInfo := new(models.AppInfo)
	appInfo.Id, err = models.GenerateAppInfoId()
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	appInfo.Appid = reqData.Appid
	appInfo.AppName = reqData.AppName
	appInfo.AppSchema = reqData.AppSchema
	appInfo.AppHost = reqData.AppHost

	appInfo.ShareClickMoney = int(reqData.ShareClickMoney)
	appInfo.ShareInstallMoney = int(reqData.ShareInstallMoney)
	appInfo.InstallMoney = int(reqData.InstallMoney)

	appInfo.Status = reqData.Status
	appInfo.Userid = userid
	appInfo.AppIcon = reqData.AppIcon
	appInfo.CreatedUTC = utils.GetNowSecond()

	err = models.InsertDBModel(nil, appInfo)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	err = caches.SetAppInfo(appInfo)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	Success(c, nil)
}

type UpdateAppPostData struct {
	Id                int64   `json:"id" binding:"required"`
	Appid             string  `json:"appid"`
	AppName           string  `json:"appname"`
	AppSchema         string  `json:"appschema"`
	AppHost           string  `json:"apphost"`
	Status            int     `json:"status"`
	ShareClickMoney   float64 `json:"share_click_money"`
	ShareInstallMoney float64 `json:"share_install_money"`
	InstallMoney      float64 `json:"install_money"`
	AppIcon           string  `json:"appicon"`
}

func UpdateApp(c *gin.Context) {
	var reqData UpdateAppPostData
	err := c.BindJSON(&reqData)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	appInfo, err := models.GetAppInfoById(nil, reqData.Id)
	if err != nil {
		Error(c, DATA_NOT_FOUND, nil, err.Error())
		return
	}

	appInfo.Appid = reqData.Appid
	appInfo.AppName = reqData.AppName
	appInfo.AppSchema = reqData.AppSchema
	appInfo.AppHost = reqData.AppHost

	appInfo.ShareClickMoney = int(reqData.ShareClickMoney)
	appInfo.ShareInstallMoney = int(reqData.ShareInstallMoney)
	appInfo.InstallMoney = int(reqData.InstallMoney)

	appInfo.Status = reqData.Status
	appInfo.AppIcon = reqData.AppIcon

	err = models.UpdateDBModel(nil, appInfo)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	err = caches.SetAppInfo(appInfo)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	Success(c, nil)
}

func AppInfoAll(c *gin.Context) {
	userid := getUserIdFromContext(c)
	if userid <= 0 {
		Error(c, LOGIN_NEEDED, nil, nil)
	}

	var res []*models.AppInfo
	//var res_ret []*models.AppInfoRet
	user, _ := models.GetUserInfoById(nil, userid)
	res, total, err := models.GetAppInfoListByUserid(nil, userid)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	if user.Name == "admin" {
		res, total, _ = models.GetAppInfoAll(nil)
		resLen := len(res)
		for i := 0; i < resLen; i++ {
			tmp_app := res[i]
			tmp_userid := tmp_app.Userid
			tmp_user, _ := models.GetUserInfoById(nil, tmp_userid)
			if tmp_user != nil {
				res[i].Des = tmp_user.Name
			}

			//temp_ret := models.GetAppInfoRet(res[i])
			//res_ret = append(res_ret, temp_ret)
		}
	} else {
		tmp_user, _ := models.GetUserInfoById(nil, userid)
		if tmp_user != nil {
			resLen := len(res)
			for i := 0; i < resLen; i++ {
				res[i].Des = tmp_user.Name
				//temp_ret := models.GetAppInfoRet(res[i])
				//res_ret = append(res_ret, temp_ret)
			}
		}
	}

	ret := gin.H{"status": true}
	ret["total"] = total
	ret["data"] = res
	c.JSON(200, ret)
}
