package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/appwilldev/sharetrace/conf"
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

	// if exist shareurl in cache, return
	var shareid int64
	shareid = 0
	old_idStr, err := caches.GetShareURLIdByUrl(postData.ShareURL)
	if err == nil && old_idStr != "" {
		log.Println("Exist cahche:%s", postData.ShareURL)
		return
	} else {
		if old_idStr != "" {
			shareid, err = strconv.ParseInt(old_idStr, 10, 64)
			if err != nil {
				Error(c, SERVER_ERROR, nil, err.Error())
				return
			}
		}
	}

	// if no shareid but exist share_url in db, reset cache, then return
	if shareid == 0 {
		old_data, err := models.GetShareURLByUrl(nil, postData.ShareURL)
		if err == nil && old_data != nil {
			//log.Println("Exist db data share_url:%s", postData.ShareURL)
			// recache
			err := caches.SetShareURL(old_data)
			if err != nil {
				log.Println("Failed to cache info %s", err.Error())
			}
			return
		}
	}

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
	err = caches.SetShareURL(data)

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

// Just return nothing, maybe  set cookie
func WebBeacon(c *gin.Context) {
	// if no share_url para, return
	q := c.Request.URL.Query()
	share_url := q["share_url"][0]
	if share_url == "" {
		log.Println("Error, No share_url para!")
		return
	}
	u, err := url.Parse(share_url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	idShareStr, err := caches.GetShareURLIdByUrl(share_url)
	var shareid int64
	shareid = 0
	if err != nil {
		log.Println(err.Error())
		// not return when redis_nil_value, for domain trace
	} else {
		shareid, err = strconv.ParseInt(idShareStr, 10, 64)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	// if no clientIP, return
	clientIP := c.ClientIP()
	if clientIP == "" {
		log.Println("No client IP")
		return
	}

	// set click_type
	click_type := conf.CLICK_TYPE_COOKIE
	agent := c.Request.Header.Get("User-Agent")
	if agent == "" {
		log.Println("No client agent")
		return
	} else {
		if strings.Contains(agent, "Safari") && !strings.Contains(agent, "Chrome") {
			click_type = conf.CLICK_TYPE_COOKIE
		} else {
			click_type = conf.CLICK_TYPE_IP
		}
	}

	md5Ctx := md5.New()
	agent_info := fmt.Sprintf("%s_%s_%s", share_url, clientIP, agent)
	md5Ctx.Write([]byte(agent_info))
	agentid := hex.EncodeToString(md5Ctx.Sum(nil))

	//log.Println("agentid:", agentid)
	_, err = caches.GetClickSessionIdByAgentId(agentid)
	if err != nil {
		log.Println("No such agentid, need create new clicksession:", agentid)
	} else {
		log.Println("Exist agentid")
		// if exist stcookieid, return
		stagentid_cookie, stagentid_err := c.Request.Cookie("stagentid")
		if stagentid_err == nil {
			if stagentid_cookie == nil || stagentid_cookie.Value == "" {
			} else {
				//log.Println("Exist stagentid:", stagentid_cookie.Value)
				return
			}
		} else {
			// need reset cookie and overwrite older agentid for buton click
			old_data, _ := models.GetClickSessionByAgentId(nil, agentid)
			log.Println("get old_data by agentid:", old_data)
			cookie := new(http.Cookie)
			if click_type == conf.CLICK_TYPE_COOKIE && old_data.Cookieid != "" {
				cookie.Name = "stcookieid"
				cookie.Value = old_data.Cookieid
				cookie.Path = "/"
				http.SetCookie(c.Writer, cookie)
			}
			if old_data.AgentId != "" {
				cookie.Name = "stagentid"
				cookie.Value = old_data.AgentId
				cookie.Path = "/"
				http.SetCookie(c.Writer, cookie)
			}
			return
		}
		return
	}

	if click_type == conf.CLICK_TYPE_COOKIE {
		// if exist stcookieid, return
		//stcookieid_cookie, stcookieid_err := c.Request.Cookie("stcookieid")
		//if stcookieid_err == nil {
		//	if stcookieid_cookie == nil || stcookieid_cookie.Value == "" {
		//	} else {
		//		//log.Println("Exist stcookieid:", stcookieid_cookie.Value)
		//		return
		//	}
		//}
	} else if click_type == conf.CLICK_TYPE_IP {
		// if cookie forbidden by client, we can use IP
		//idStr, err := caches.GetClickSessionIdByIP(clientIP)
		//// if already exist IP cache, recookie, return
		//if err == nil && idStr != "" {
		//	log.Println("Exist IP:", clientIP)
		//	// TODO:没有 cookie，但是同IP， 不同Agent呢？
		//	// if exist stagentid, return
		//	_, err = caches.GetClickSessionIdByAgentId(agentid)
		//	if err != nil {
		//		log.Println("Exist IP but no such agentid:", agentid)
		//	} else {
		//		return
		//	}
		//}
	}

	id, err := models.GenerateClickSessionId()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// insert to db
	data := new(models.ClickSession)
	data.Id = id
	data.ClickType = click_type
	data.Agent = agent
	data.AgentIP = clientIP
	data.CreatedUTC = utils.GetNowSecond()
	data.URLHost = u.Host
	data.ClickURL = share_url

	if shareid > 0 {
		data.Shareid = shareid
		data.Cookieid = fmt.Sprintf("%s_%d_%d", conf.COOKIE_PREFIX, shareid, id)
	}

	data.AgentId = agentid

	log.Println("Webbeacon clicksession data:%s", data)

	err = models.InsertDBModel(nil, data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// cache clicksession by cookie / IP
	err = caches.SetClickSession(data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	cookie := new(http.Cookie)
	if click_type == conf.CLICK_TYPE_COOKIE && data.Cookieid != "" {
		cookie.Name = "stcookieid"
		//cookie.Expires = time.Now().Add(time.Duration(7*86400) * time.Second)
		cookie.Value = data.Cookieid
		cookie.Path = "/"
		http.SetCookie(c.Writer, cookie)
	}

	if data.AgentId != "" {
		cookie.Name = "stagentid"
		cookie.Value = data.AgentId
		cookie.Path = "/"
		http.SetCookie(c.Writer, cookie)
	}

	return
}

func ClickInstallButton(c *gin.Context) {
	agentid := ""
	old_cookie, err := c.Request.Cookie("stagentid")
	if err == nil && old_cookie != nil && old_cookie.Value != "" {
		agentid = old_cookie.Value
	} else {
		log.Println("err:", err)
		return
	}

	q := c.Request.URL.Query()
	buttonid := q["buttonid"][0]
	if buttonid == "" || buttonid == "undefined" {
		log.Println("No buttonid para:")
	} else {
		log.Println("buttonid:", buttonid)
	}

	idStr, err := caches.GetClickSessionIdByAgentId(agentid)
	if err != nil {
		log.Println("No cache or db! Data:", agentid)
		return
	} else {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		data, _ := caches.GetClickSessionModelInfoById(id)

		if data.Status == conf.CLICK_SESSION_STATUS_CLICK {
			data.Status = conf.CLICK_SESSION_STATUS_BUTTON
			if buttonid != "" {
				data.ButtonId = buttonid
			}
			err = models.UpdateDBModel(nil, data)
			if err != nil {
				Error(c, SERVER_ERROR, nil, err.Error())
				return
			}
			err = caches.UpdateClickSession(data)
			if err != nil {
				Error(c, SERVER_ERROR, nil, err.Error())
				return
			}
		} else {
			log.Println("Duplicated clickbutton! Data:", data)
		}

	}

}

func WebBeaconCheck(c *gin.Context) {
	q := c.Request.URL.Query()
	appid := q["appid"][0]
	if appid == "" {
		Error(c, BAD_REQUEST, nil, "No Appid")
		return
	}

	app, err := models.GetAppInfoByAppid(nil, appid)
	if err != nil || app == nil {
		if err == nil {
			err = fmt.Errorf("app nil")
		}
		Error(c, DATA_NOT_FOUND, nil, err.Error())
		return
	}

	installid := q["installid"][0]
	if installid == "" {
		Error(c, BAD_REQUEST, nil, "No installid")
		return
	}

	appschema := app.AppSchema

	// TODO get cookie, ip, return , then  write to db
	idStr := ""
	old_cookie, err := c.Request.Cookie("stcookieid")
	click_type := conf.CLICK_TYPE_COOKIE
	trackid := ""

	//cookie := new(http.Cookie)
	//cookie.Name = "stcookieid"
	//cookie.Expires = time.Now().Add(-3600)
	//cookie.Value = ""
	//cookie.Path = "/"
	//http.SetCookie(c.Writer, cookie)

	data := gin.H{"appschema": appschema}
	c.HTML(200, "webbeaconcheck.html", data)

	if err == nil && old_cookie != nil && old_cookie.Value != "" {
		click_type = conf.CLICK_TYPE_COOKIE
		trackid = old_cookie.Value
	} else {
		click_type = conf.CLICK_TYPE_IP
		trackid = c.ClientIP()
	}

	idStr, _ = caches.GetClickSessionId(click_type, trackid)
	if idStr == "" {
		old_data, err := models.GetClickSession(nil, click_type, trackid)
		if err == nil && old_data != nil {
			log.Println("No cache for data and recached")
			_ = caches.SetClickSession(old_data)
			idStr, _ = caches.GetClickSessionId(click_type, trackid)
		} else {
			Error(c, SERVER_ERROR, nil, err.Error())
			return
		}
	}

	if idStr == "" {
		log.Println("No cache or db! Data:", click_type, trackid)
		// TODO newuser not from share
		return
	} else {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		data, _ := caches.GetClickSessionModelInfoById(id)

		if data.Status == conf.CLICK_SESSION_STATUS_CLICK {
			data.Installid = installid
			data.ClickType = click_type
			data.Status = conf.CLICK_SESSION_STATUS_INSTALLED
			err = models.UpdateDBModel(nil, data)
			if err != nil {
				Error(c, SERVER_ERROR, nil, err.Error())
				return
			}
			err = caches.UpdateClickSession(data)
			if err != nil {
				Error(c, SERVER_ERROR, nil, err.Error())
				return
			}
		} else {
			log.Println("Duplicated webbeaconcheck! Data:", trackid)
		}
	}
}
