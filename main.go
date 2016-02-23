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
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

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

	stAPIV1 := ginIns.Group("/1/")
	{
		stAPIV1.POST("/st/share", controllers.Share)
		stAPIV1.POST("/st/click", controllers.Click)
		stAPIV1.POST("/st/install", controllers.Install)
		stAPIV1.GET("/st/score", controllers.Score)
		stAPIV1.GET("/st/webbeacon", controllers.WebBeacon)
	}

	// op api
	//opAPIGroup := ginIns.Group("/op")
	//{
	//	//opAPIGroup.POST("/login", Login)
	//	//opAPIGroup.POST("/logout", OpAuth, Logout)

	//	//opAPIGroup.GET("/users/:page/:count", InitUserCheck, OpAuth, GetUsers)
	//	//opAPIGroup.POST("/user", OpAuth, ConfWriteCheck, NewUser)
	//	//opAPIGroup.PUT("/user", OpAuth, ConfWriteCheck, UpdateUser)
	//	//opAPIGroup.POST("/user/init", ConfWriteCheck, InitUser)
	//	//opAPIGroup.GET("/user/info", OpAuth, GetLoginUserInfo)

	//	//opAPIGroup.GET("/apps/user/:user_key", OpAuth, GetApps)
	//	//opAPIGroup.GET("/apps/all/:page/:count", OpAuth, GetAllApps)
	//	//opAPIGroup.GET("/app/:app_key", OpAuth, GetApp)
	//	//opAPIGroup.GET("/apps/search", OpAuth, SearchApps)
	//	//opAPIGroup.POST("/app", OpAuth, ConfWriteCheck, NewApp)
	//	//opAPIGroup.PUT("/app", OpAuth, ConfWriteCheck, UpdateApp)

	//}

	gracehttp.Serve(&http.Server{Addr: conf.HttpAddr, Handler: ginIns})
}
