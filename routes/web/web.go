package web

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func editServerInfo(r *gin.Engine, handle *handle.WebHandler) {
	r.POST("/web/serverinfo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.EditServerInfo(c),
		})
	})
}

func getServerInfo(r *gin.Engine, handle *handle.WebHandler) {
	r.GET("/web/serverinfo", func(c *gin.Context) {
		data, code := handle.GetServerInfo()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func verifyFtAccount(r *gin.Engine, handle *handle.WebHandler) {
	r.POST("/web/verifyftaccount", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.VerifyFtAccount(c),
		})
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.WebHandler) {
	editServerInfo(r, handle)
	getServerInfo(r, handle)
	verifyFtAccount(r, handle)
}
