package main

import (
	// "app/event"

	// "app/tools"

	"net/http"
	"runtime"
	"shensuanzi/commondata"
	"shensuanzi/conf"
	"shensuanzi/datastruct"
	"shensuanzi/handle"
	"shensuanzi/routes/app"
	"shensuanzi/routes/web"

	"github.com/facebookgo/grace/gracehttp"

	"github.com/gin-gonic/gin"
)

//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Appversion, Apptoken, IsPc,Platform,Webtoken")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		c.Next()
	}
}

func checkVersion(handle *handle.AppHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := "/web/serverinfo"
		if c.Request.Method == "POST" && c.Request.RequestURI == url {
			c.Next()
			return
		}
		serverVersion, isMaintain := handle.GetServerInfoFromMemory()
		if isMaintain == true {
			c.AbortWithStatus(int(datastruct.Maintenance))
			return
		}
		version, isExist := c.Request.Header["Appversion"]
		if isExist && version[0] == serverVersion {
			c.Next()
		} else {
			c.AbortWithStatusJSON(int(datastruct.VersionError), handle.GetDirectDownloadApp())
		}
	}
}

func createData() (*handle.AppHandler, *handle.WebHandler) {
	app_hanle := handle.CreateAppHandle()
	web_hanle := handle.CreateWebHandle()
	commondata.Create()
	serverInfo := handle.CreateServerInfo(web_hanle.GetWebDBHandler())
	commondata.SetServerInfo(serverInfo)
	return app_hanle, web_hanle
}

func registerRoutes(r *gin.Engine, app_hanle *handle.AppHandler, web_hanle *handle.WebHandler) {
	app.FtRegisterRoutes(r, app_hanle)
	app.UserRegisterRoutes(r, app_hanle)
	web.RegisterRoutes(r, web_hanle)
}

func main() {
	r := gin.Default()
	var mode string
	switch conf.Common.Mode {
	case conf.Debug:
		mode = gin.DebugMode
	case conf.Release:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r.Use(cors())

	app_hanle, web_hanle := createData()
	r.Use(checkVersion(app_hanle))

	registerRoutes(r, app_hanle, web_hanle)

	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "linux":
		server := &http.Server{Addr: conf.Server.HttpServer, Handler: r}
		gracehttp.Serve(server)
	case "windows":
		r.Run(conf.Server.HttpServer)
	}

	// r.Run(conf.Server.HttpServer) //listen and serve on 0.0.0.0:8080 asd
}
