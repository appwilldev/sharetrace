package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/appwilldev/sharetrace/models"
	"github.com/appwilldev/sharetrace/models/caches"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/gin-gonic/gin"
)

func Share(c *gin.Context) {
	var postData struct {
		ShareURL string `json:"share_url" binding:"required"`
		Fromid   string `json:"fromid" binding:"required"`
		Appid    string `json:"appid" binding:"required"`
		Itemid   string `json:"itemid"`
		Channel  string `json:"channel"`
		Ver      string `json:"ver"`
		Des      string `json:"des"`
	}
	err := c.BindJSON(&postData)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	log.Println("Share data:%s", postData)

	data := new(models.ShareURL)
	id, err := models.GenerateShareURLId()
	data.Id = id
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	data.ShareURL = postData.ShareURL
	data.Fromid = postData.Fromid
	data.Appid = postData.Appid

	if postData.Itemid != "" {
		data.Itemid = postData.Itemid
	}
	if postData.Channel != "" {
		data.Channel = postData.Channel
	}
	if postData.Ver != "" {
		data.Ver = postData.Ver
	}
	if postData.Des != "" {
		data.Des = postData.Des
	}

	data.CreatedUTC = utils.GetNowSecond()

	log.Println("Generate data:%s", data)

	err = models.InsertDBModel(nil, data)
	if err != nil {
		Error(c, SERVER_ERROR, nil, nil)
		return
	}
	err = caches.NewShareURL(data)

	ret := gin.H{"status": true}
	c.JSON(200, ret)
}

func Click(c *gin.Context) {
	var postData struct {
		ShareURL string `json:"share_url" binding:"required"`
	}
	err := c.BindJSON(&postData)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	log.Println("Click data:%s", postData)

	idStr, err := caches.GetShareURLIdByUrl(postData.ShareURL)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	log.Println("idStr:", idStr)
	shareid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	log.Println("Shareid:", shareid)

	data := new(models.ClickSession)
	id, err := models.GenerateClickSessionId()
	log.Println("Clickid:", id)
	data.Id = id
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	data.Shareid = shareid
	// generate cookieid
	cookieid := fmt.Sprintf("st_%d_%d", shareid, id)
	data.Cookieid = cookieid

	data.CreatedUTC = utils.GetNowSecond()

	log.Println("Generate data:%s", data)

	err = models.InsertDBModel(nil, data)
	if err != nil {
		Error(c, SERVER_ERROR, nil, nil)
		return
	}
	err = caches.NewClickSession(data)

	ret := gin.H{"status": true}
	ret["st_cookieid"] = cookieid
	c.JSON(200, ret)
}

func AgentClick(c *gin.Context) {
	var postData struct {
		ShareURL string `json:"share_url" binding:"required"`
		AgentIP  string `json:"agent_ip" binding:"required"`
		Agent    string `json:"agent"`
	}
	err := c.BindJSON(&postData)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	log.Println("Agent Click data:%s", postData)

	idStr, err := caches.GetShareURLIdByUrl(postData.ShareURL)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	log.Println("idStr:", idStr)
	shareid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	log.Println("Shareid:", shareid)

	data := new(models.ClickSession)
	id, err := models.GenerateClickSessionId()
	log.Println("Clickid:", id)
	data.Id = id
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	data.Shareid = shareid

	cookieid := fmt.Sprintf("st_%d_%d", shareid, id)
	data.Cookieid = cookieid

	data.ClickType = 1
	data.Agent = postData.Agent
	data.AgentIP = postData.AgentIP
	// todo: generate agentid

	data.CreatedUTC = utils.GetNowSecond()

	log.Println("Generate data:%s", data)

	err = models.InsertDBModel(nil, data)
	if err != nil {
		Error(c, SERVER_ERROR, nil, nil)
		return
	}
	err = caches.NewClickSession(data)

	ret := gin.H{"status": true}
	ret["st_cookieid"] = cookieid
	c.JSON(200, ret)
}

func Install(c *gin.Context) {
	var postData struct {
		St_cookieid string `json:"st_cookieid" binding:"required"`
		Installid   string `json:"installid" binding:"required"`
	}
	err := c.BindJSON(&postData)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	log.Println("Install data:%s", postData)

	idStr, err := caches.GetClickSessionIdByCookieid(postData.St_cookieid)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	data, err := caches.GetClickSessionModelInfoById(id)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	if data.Installid == "" {
		data.Installid = postData.Installid
		err = models.UpdateDBModel(nil, data)
		if err != nil {
			Error(c, SERVER_ERROR, nil, err.Error())
			//return
		}
		err = caches.UpdateClickSession(data)
		if err != nil {
			Error(c, SERVER_ERROR, nil, err.Error())
			//return
		}

	} else {
		ret := gin.H{"status": true}
		c.JSON(200, ret)
		return
	}
	ret := gin.H{"status": true}
	c.JSON(200, ret)
}

func Score(c *gin.Context) {
	q := c.Request.URL.Query()
	userIdStr := q["userid"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	log.Println("Score userId:", userId)

	appIdStr := q["appid"][0]
	appId, _ := strconv.ParseInt(appIdStr, 10, 64)
	log.Println("Score appId:", appId)

	dataList, _ := models.GetShareClickListOfAppUser(nil, appIdStr, userIdStr)
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
	//c.HTML(200, "mymoney.html", data)
	return
}

func WebBeacon(c *gin.Context) {
	q := c.Request.URL.Query()
	share_url := q["share_url"][0]
	if share_url == "" {
		Error(c, DATA_NOT_FOUND, nil, "No URL")
		return
	}
	log.Println("share_url:", share_url)

	clientIP := c.ClientIP()
	if clientIP == "" {
		Error(c, SERVER_ERROR, nil, "No ClientIP")
		return
	}

	click_type := 0
	agent := c.Request.Header.Get("User-Agent")
	if agent == "" {
		Error(c, DATA_NOT_FOUND, nil, "No Agent")
		return
	} else {
		if strings.Contains(agent, "Safari") {
			click_type = 0
		} else {
			click_type = 1
		}
	}

	old_cookie, err := c.Request.Cookie("stcookieid")
	if err == nil {
		log.Println("Exist stcookieid:", old_cookie.Value)
		return
	}

	idStr, err := caches.GetShareURLIdByUrl(share_url)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}
	shareid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	data := new(models.ClickSession)
	id, err := models.GenerateClickSessionId()
	data.Id = id
	if err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	data.Shareid = shareid
	cookieid := fmt.Sprintf("st_%d_%d", shareid, id)
	data.Cookieid = cookieid
	data.ClickType = click_type
	data.Agent = agent
	data.AgentIP = clientIP

	data.CreatedUTC = utils.GetNowSecond()

	log.Println("Generate data:%s", data)

	err = models.InsertDBModel(nil, data)
	if err != nil {
		Error(c, SERVER_ERROR, nil, nil)
		return
	}
	err = caches.NewClickSession(data)

	cookie := new(http.Cookie)
	cookie.Name = "stcookieid"
	cookie.Expires = time.Now().Add(time.Duration(7*86400) * time.Second)
	cookie.Value = cookieid
	cookie.Path = "/"
	http.SetCookie(c.Writer, cookie)
	return
}

func WebBeaconCheck(c *gin.Context) {
	q := c.Request.URL.Query()
	appid := q["appid"][0]
	if appid == "" {
		Error(c, DATA_NOT_FOUND, nil, "No Appid")
		return
	}
	log.Println("appid:", appid)

	// TODO get appschema by appid

	data := gin.H{"appschema": "avft"}
	c.HTML(200, "webbeaconcheck.html", data)
}
