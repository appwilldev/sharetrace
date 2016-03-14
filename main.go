package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	//"strings"

	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/controllers"
	"github.com/appwilldev/sharetrace/logger"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		dur := time.Since(start) / time.Millisecond

		requestLogData := controllers.GetRequestLogDataFromContext(c)
		logInfo := map[string]interface{}{
			"code":   c.Writer.Status(),
			"dur":    dur,
			"remote": c.ClientIP(),
			"url":    c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"method": c.Request.Method,
			"data":   requestLogData,
		}

		if c.Writer.Status() >= 400 {
			logger.RequestLogger.Error(logInfo)
		} else {
			logger.RequestLogger.Info(logInfo)
		}
	}
}

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
	ginIns.Use(requestLogger())

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

	stAPIV1 := ginIns.Group("/1/st")
	{
		stAPIV1.POST("/share", authHandler(), controllers.Share)
		stAPIV1.GET("/webbeacon", controllers.WebBeacon)
		stAPIV1.GET("/webbeaconbutton", controllers.ClickInstallButton)
		stAPIV1.GET("/webbeaconcheck", controllers.WebBeaconCheck)
		stAPIV1.GET("/score", controllers.Score)
	}

	statsAPIV1 := ginIns.Group("/1/stats")
	{
		statsAPIV1.GET("/share", controllers.StatsShare)
		statsAPIV1.GET("/total", controllers.StatsTotal)
		statsAPIV1.GET("/host", controllers.StatsHost)
	}

	//ctAPIV1 := ginIns.Group("/1/ct")
	//{
	//	ctAPIV1.GET("/webbeacon", controllers.WebBeaconCT)
	//}

	// op api
	opAPIGroup := ginIns.Group("/op")
	{
		// user
		opAPIGroup.POST("/user/init", controllers.Register)
		opAPIGroup.POST("/login", controllers.Login)
		opAPIGroup.POST("/logout", authHandler(), controllers.Logout)
		opAPIGroup.GET("/users/:page/:count", authHandler(), controllers.UserInfoAll)

		// app
		opAPIGroup.GET("/apps/all/:page/:count", authHandler(), controllers.AppInfoAll)
		opAPIGroup.POST("/app", authHandler(), controllers.NewApp)
		opAPIGroup.PUT("/app", authHandler(), controllers.UpdateApp)

	}

	ginIns.LoadHTMLFiles("./templates/webbeaconcheck.html", "./templates/score.html")

	gracehttp.Serve(&http.Server{Addr: conf.HttpAddr, Handler: ginIns})
}
