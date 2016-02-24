package controllers

import (
	"log"
	//"net/http"
	//"strconv"

	"github.com/appwilldev/sharetrace/models"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/gin-gonic/gin"
)

type NewAppPostData struct {
	Appid   string `json:"appid" binding:"required"`
	AppName string `json:"appname" binding:"required"`
	AppIcon string `json:"appicon"`
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
	appInfo.Userid = userid
	appInfo.AppIcon = reqData.AppIcon
	appInfo.CreatedUTC = utils.GetNowSecond()

	err = models.InsertDBModel(nil, appInfo)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	Success(c, nil)
}

func AppInfoAll(c *gin.Context) {
	var res interface{}
	res, total, _ := models.GetAppInfoAll(nil)
	ret := gin.H{"status": true}
	ret["total"] = total
	ret["data"] = res
	c.JSON(200, ret)
}
