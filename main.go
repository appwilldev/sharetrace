package main

import (
	"log"
	"net/http"

	//"time"

	"os"
	"path/filepath"
	"strconv"

	//"strings"

	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/controllers"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

func authHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(utils.CookieKey)
		if err != nil {
			return
		}

		userId := utils.DecodeCookie(cookie.Value)
		if conf.DebugMode {
			log.Println("get user:", userId)
		}
		if userId > 0 {
			c.Set("userid", userId)
		}

		c.Next()
	}
}

func main() {
	wd, _ := os.Getwd()
	pidFile, err := os.OpenFile(filepath.Join(wd, "sharetrace.pid"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("failed to create pid file: %s", err.Error())
		os.Exit(1)
	}
	pidFile.WriteString(strconv.Itoa(os.Getpid()))
	pidFile.Close()

	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginIns := gin.New()
	ginIns.Use(gin.Recovery())

	if conf.DebugMode {
		ginIns.Use(gin.Logger())
	}

	if conf.WebDebugMode {
		// static
		ginIns.Static("/web", "./web")
	} else {
		// bin static
		ginIns.GET("/web/*file",
			func(c *gin.Context) {
				fileName := c.Param("file")
				if fileName == "/" {
					fileName = "/index.html"
				}
				//TOASK : ./main.go:53: undefined: Asset
				//data, err := Asset("web" + fileName)
				//if err != nil {
				//	c.String(http.StatusNotFound, err.Error())
				//	return
				//}

				//switch {
				//case strings.LastIndex(fileName, ".html") == len(fileName)-5:
				//	c.Header("Content-Type", "text/html; charset=utf-8")
				//case strings.LastIndex(fileName, ".css") == len(fileName)-4:
				//	c.Header("Content-Type", "text/css")
				//}
				//c.String(http.StatusOK, string(data))
			})
	}

	userAPIV1 := ginIns.Group("/1/user")
	{
		userAPIV1.POST("/register", controllers.Register)
		userAPIV1.POST("/login", controllers.Login)
		userAPIV1.POST("/logout", authHandler(), controllers.Logout)
		userAPIV1.GET("/all", controllers.UserInfoAll)
	}

	appAPIV1 := ginIns.Group("/1/app")
	{
		appAPIV1.POST("/new", authHandler(), controllers.NewApp)
		appAPIV1.PUT("/update", authHandler(), controllers.UpdateApp)
		appAPIV1.GET("/all", controllers.AppInfoAll)
	}

	statsAPIV1 := ginIns.Group("/1/stats")
	{
		statsAPIV1.GET("/share", controllers.StatsShare)
		statsAPIV1.GET("/total", controllers.StatsTotal)
	}

	stAPIV1 := ginIns.Group("/1/st")
	{
		stAPIV1.POST("/share", authHandler(), controllers.Share)
		stAPIV1.POST("/click", authHandler(), controllers.Click)
		stAPIV1.POST("/agentclick", authHandler(), controllers.AgentClick)
		stAPIV1.POST("/install", authHandler(), controllers.Install)
		stAPIV1.GET("/score", controllers.Score)
		stAPIV1.GET("/webbeacon", controllers.WebBeacon)
		stAPIV1.GET("/webbeaconcheck", controllers.WebBeaconCheck)
	}

	// op api
	opAPIGroup := ginIns.Group("/op")
	{
		opAPIGroup.POST("/user/init", controllers.Register)
		opAPIGroup.POST("/login", controllers.Login)
		opAPIGroup.POST("/logout", authHandler(), controllers.Logout)

		opAPIGroup.GET("/users/:page/:count", authHandler(), controllers.UserInfoAll)
		//opAPIGroup.POST("/user", OpAuth, ConfWriteCheck, NewUser)
		//opAPIGroup.PUT("/user", OpAuth, ConfWriteCheck, UpdateUser)
		//opAPIGroup.GET("/user/info", OpAuth, GetLoginUserInfo)

		//opAPIGroup.GET("/apps/user/:user_key", OpAuth, GetApps)
		opAPIGroup.GET("/apps/all/:page/:count", controllers.AppInfoAll)
		opAPIGroup.POST("/app", authHandler(), controllers.NewApp)
		opAPIGroup.PUT("/app", authHandler(), controllers.UpdateApp)

	}

	ginIns.LoadHTMLFiles("./templates/webbeaconcheck.html")

	gracehttp.Serve(&http.Server{Addr: conf.HttpAddr, Handler: ginIns})
}
