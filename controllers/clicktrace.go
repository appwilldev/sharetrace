package controllers

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/url"
	//"strconv"
	//"strings"
	//"time"
	"encoding/hex"

	//"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
	//"github.com/appwilldev/sharetrace/models/caches"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/gin-gonic/gin"
)

// Just return nothing, maybe  set cookie
// ClickTrace WebBeacon
func WebBeaconCT(c *gin.Context) {

	// if exist stcookieid, return
	old_cookie, err := c.Request.Cookie("stagentid")
	if err == nil {
		if old_cookie == nil || old_cookie.Value == "" {
		} else {
			log.Println("Exist stagentid:", old_cookie.Value)
			return
		}
	}

	// if no share_url para, return
	q := c.Request.URL.Query()
	share_url := q["share_url"][0]
	if share_url == "" {
		log.Println("No share_url para:")
		return
	}
	u, _ := url.Parse(share_url)
	log.Println("url pares host:", u.Host)
	//log.Println("share_url:", share_url)

	// if no clientIP, return
	clientIP := c.ClientIP()
	if clientIP == "" {
		log.Println("No client IP")
		return
	}

	agent := c.Request.Header.Get("User-Agent")
	if agent == "" {
		log.Println("No client agent")
		return
	}

	data := new(models.ClickTrace)
	id, err := models.GenerateClickTraceId()
	data.Id = id
	if err != nil {
		log.Println(err.Error())
		return
	}

	data.Agent = agent
	data.AgentIP = clientIP
	data.URLHost = u.Host

	md5Ctx := md5.New()
	agent_info := fmt.Sprintf("%s_%s_%s", share_url, clientIP, agent)
	//log.Println("agent_info:", agent_info)
	md5Ctx.Write([]byte(agent_info))
	cipherStr := hex.EncodeToString(md5Ctx.Sum(nil))
	log.Println(cipherStr)

	agentId := cipherStr
	data.AgentId = agentId

	data.CreatedUTC = utils.GetNowSecond()

	log.Println("Webbeacon clicktrace data:%s", data)

	// insert to db
	err = models.InsertDBModel(nil, data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "stagentid"
	cookie.Value = agentId
	cookie.Path = "/"
	http.SetCookie(c.Writer, cookie)

	return
}
