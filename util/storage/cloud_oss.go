package storage

import (
	"github.com/casdoor/oss"
	"github.com/casdoor/oss/aliyun"
	"github.com/casdoor/oss/azureblob"
)

func NewAliyunOssStorageProvider(clientId string, clientSecret string, bucket string, endpoint string) oss.StorageInterface {
	sp := aliyun.New(&aliyun.Config{
		AccessID:  clientId,
		AccessKey: clientSecret,
		Bucket:    bucket,
		Endpoint:  endpoint,
	})

	return sp
}

func NewAzureBlobStorageProvider(clientId string, clientSecret string, region string, bucket string, endpoint string) oss.StorageInterface {
	sp := azureblob.New(&azureblob.Config{
		AccessId:  clientId,
		AccessKey: clientSecret,
		Region:    region,
		Bucket:    bucket,
		Endpoint:  endpoint,
	})
	return sp
}
