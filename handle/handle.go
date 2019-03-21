package handle

import (
	"log"
	"shensuanzi/cache"
	"shensuanzi/datastruct"
	"shensuanzi/db"
	"sync"
)

type AppHandler struct {
	dbHandler    *db.DBHandler
	cacheHandler *cache.CACHEHandler
	payMutex     *sync.Mutex //读写互斥量
}

var appHandler *AppHandler

func CreateAppHandle() *AppHandler {
	appHandler = new(AppHandler)
	appHandler.cacheHandler = cache.CreateCACHEHandler()
	appHandler.dbHandler = db.CreateDBHandler(false)
	return appHandler
}

type WebHandler struct {
	dbHandler *db.DBHandler
}

var webHandler *WebHandler

func CreateWebHandle() *WebHandler {
	webHandler = new(WebHandler)
	webHandler.dbHandler = db.CreateDBHandler(true)
	return webHandler
}

func (handler *WebHandler) GetWebDBHandler() *db.DBHandler {
	return handler.dbHandler
}

func CreateServerInfo(dbHandler *db.DBHandler) *datastruct.ServerData {
	server := new(datastruct.ServerData)
	server.RWMutex = new(sync.RWMutex)
	db_data, code := dbHandler.GetServerInfo()
	if code != datastruct.NULLError {
		log.Fatal("GetServerInfo error from db")
		return nil
	}
	server.IsMaintain = db_data.IsMaintain
	server.Version = db_data.Version
	return server
}
