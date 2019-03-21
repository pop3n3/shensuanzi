package datastruct

type WebServerInfoBody struct {
	Version    string `json:"version"`
	IsMaintain bool   `json:"ismaintain"`
	GzhAppid   string `json:"gzh_appid"`
	KfptAppid  string `json:"kfpt_appid"`
}

type WebVerifyFtAccountBody struct {
	IsPassed bool `json:"ispassed"`
	FtId     int  `json:"ftid"`
}
