package handle

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct"

	"github.com/gin-gonic/gin"
)

func (web *WebHandler) EditServerInfo(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebServerInfoBody
	err := c.BindJSON(&body)
	if err != nil || body.GzhAppid == "" || body.KfptAppid == "" || body.Version == "" {
		return datastruct.ParamError
	}
	code := web.dbHandler.EditServerInfo(&body)
	if code != datastruct.NULLError {
		return code
	}
	serverInfo := commondata.GetServerInfo()
	serverInfo.RWMutex.Lock()
	defer serverInfo.RWMutex.Unlock()
	serverInfo.Version = body.Version
	serverInfo.IsMaintain = body.IsMaintain
	return datastruct.NULLError
}

func (web *WebHandler) GetServerInfo() (interface{}, datastruct.CodeType) {
	return web.dbHandler.GetServerInfo()
}

func (web *WebHandler) VerifyFtAccount(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebVerifyFtAccountBody
	err := c.BindJSON(&body)
	if err != nil || body.FtId <= 0 {
		return datastruct.ParamError
	}
	return web.dbHandler.VerifyFtAccount(&body)
}
