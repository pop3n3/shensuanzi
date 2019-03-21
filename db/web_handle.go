package db

import (
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"shensuanzi/tools"
)

func (handle *DBHandler) GetServerInfo() (*datastruct.WebServerInfoBody, datastruct.CodeType) {
	engine := handle.mysqlEngine
	serverInfo := new(datastruct.ServerInfo)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(serverInfo)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebServerInfoBody)
	resp.IsMaintain = serverInfo.IsMaintain
	resp.Version = serverInfo.Version
	resp.GzhAppid = serverInfo.GzhAppid
	resp.KfptAppid = serverInfo.KfptAppid
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditServerInfo(body *datastruct.WebServerInfoBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	serverInfo := new(datastruct.ServerInfo)
	serverInfo.GzhAppid = body.GzhAppid
	serverInfo.IsMaintain = body.IsMaintain
	serverInfo.KfptAppid = body.KfptAppid
	serverInfo.Version = body.Version

	_, err := engine.Where("id=?", datastruct.DefaultId).Cols("version", "is_maintain", "gzh_appid", "kfpt_appid").Update(serverInfo)
	if err != nil {
		log.Error("EditServerInfo err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) VerifyFtAccount(body *datastruct.WebVerifyFtAccountBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	rs := tools.BoolToAuthState(body.IsPassed)
	sql := "update cold_f_t_info set auth_state=? where id=?"
	_, err := engine.Exec(sql, rs, body.FtId)
	if err != nil {
		log.Error("VerifyFtAccount err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}
