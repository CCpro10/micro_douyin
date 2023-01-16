package util

import (
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type OSSClient struct {
	bucket *oss.Bucket
}

var ossClient *OSSClient

func GetOSSClient() *OSSClient {
	return ossClient
}

func InitOSSClient(conf *conf.Conf) {
	client, err := oss.New(conf.Oss.Endpoint, conf.Oss.AccessKeyId, conf.Oss.AccessKeySecret)
	if err != nil {
		log.Panicln("oss.New Failed:", err)
		return
	}
	bucket, err := client.Bucket(conf.Oss.BucketName)
	if err != nil {
		log.Panicln("client.Bucket() Failed:", err)
		return
	}
	ossClient = &OSSClient{bucket: bucket}
}

func (c *OSSClient) UploadFileFromStream(objectName string, data io.Reader) error {
	return c.bucket.PutObject(objectName, data)
}

func (c *OSSClient) UploadFileFromLocalFile(objectName string, filePath string) error {
	return c.bucket.PutObjectFromFile(objectName, filePath)
}
