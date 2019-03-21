package osstool

import (
	"fmt"
	"os"
	"shensuanzi/datastruct/important"
	"shensuanzi/log"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func CreateOSSBucket() *oss.Bucket {
	// 创建OSSClient实例。
	client, err := oss.New(important.OSSEndpoint, important.OSSAccessKeyId, important.OSSAccessKeySecret)
	if err != nil {
		log.Debug("Error:%v", err)
		os.Exit(-1)
	}

	bucketName := important.OSSBucketName

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Debug("Error:%v", err)
		os.Exit(-1)
	}

	return bucket
}

func DeleteFile(bucket *oss.Bucket, url string) {
	str_arr := strings.Split(url, "/")
	count := len(str_arr)
	objectName := ""
	if count >= 2 {
		objectName = str_arr[count-2] + "/" + str_arr[count-1]
	}
	if objectName == "" {
		return
	}
	err := bucket.DeleteObject(objectName)
	if err != nil {
		log.Error("osstool DeleteFile Error:%v", err)
	}
	//https://shensuanzi.oss-cn-shenzhen.aliyuncs.com/ft_avatar_dev/110485312978812928.png
}

// func SignedURL(bucket *oss.Bucket, objectName string) string {
// 	// signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, 6000)
// 	// //oss.Process("image/resize,h_100"))
// 	// if err != nil {
// 	// 	log.Debug("SignedURL err:%v", err.Error())
// 	// 	return ""
// 	// }
// 	signedURL := fmt.Sprintf("https://rouge999.oss-cn-shenzhen.aliyuncs.com/%s", objectName)
// 	return signedURL
// }

func CreateOSSURL(objectName string) string {
	url := fmt.Sprintf("https://rouge999.oss-cn-shenzhen.aliyuncs.com/%s", objectName)
	return url
}
