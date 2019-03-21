package commondata

import (
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"shensuanzi/osstool"
	"strconv"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/holdno/snowFlakeByGo"
)

type CommonData struct {
	idWorker   *snowFlakeByGo.Worker
	serverInfo *datastruct.ServerData
	ossBucket  *oss.Bucket
}

var CommonDataInfo *CommonData
var once sync.Once

func Create() *CommonData {
	once.Do(func() {
		CommonDataInfo = new(CommonData)
		idWorker, err := snowFlakeByGo.NewWorker(0)
		if err != nil {
			log.Fatal("CreateCommonData err:%v", err.Error())
		}
		CommonDataInfo.idWorker = idWorker
		CommonDataInfo.ossBucket = osstool.CreateOSSBucket()
	})
	return CommonDataInfo
}

func UniqueId() string {
	return strconv.FormatInt(CommonDataInfo.idWorker.GetId(), 10)
}

func GetServerInfo() *datastruct.ServerData {
	return CommonDataInfo.serverInfo
}

func SetServerInfo(data *datastruct.ServerData) {
	CommonDataInfo.serverInfo = data
}

func DeleteOSSFileWithUrl(url string) {
	osstool.DeleteFile(CommonDataInfo.ossBucket, url)
}
