package util

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/CCpro10/micro_douyin/conf"
)

//func TestOSSUploadFileFromLocalFile(t *testing.T) {
//	confPath := "../conf/conf.yaml"
//	conf.Init(confPath)
//	InitOSSClient(conf.Config)
//
//	objectName := fmt.Sprintf("test/t0_%d.jpg", time.Now().Unix())
//	filePath := "/Users/qizong007/Downloads/get.jpg"
//
//	err := GetOSSClient().UploadFileFromLocalFile(objectName, filePath)
//	if err != nil {
//		log.Println(err)
//	}
//}

func TestOSSUploadFileFromStream(t *testing.T) {
	confPath := "../conf/conf.yaml"
	conf.Init(confPath)
	InitOSSClient(conf.Config)

	objectName := fmt.Sprintf("test/t1_%d.jpg", time.Now().Unix())
	url := "https://i0.hdslb.com/bfs/album/d89a115bd24fc3989f2c03bbb92faf64582c0d69.png"
	res, _ := http.Get(url)

	_, err := GetOSSClient().Put(objectName, res.Body)
	if err != nil {
		log.Println(err)
	}
}
