package util

import (
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/CCpro10/micro_douyin/util/storage"
	"github.com/casdoor/oss"
)

func GetStorageProvider(providerType string, clientId string, clientSecret string, region string, bucket string, endpoint string) oss.StorageInterface {
	switch providerType {
	case "Local File System":
		return storage.NewLocalFileSystemStorageProvider()
	case "Aliyun OSS":
		return storage.NewAliyunOssStorageProvider(clientId, clientSecret, bucket, endpoint)
	case "Azure Blob":
		return storage.NewAzureBlobStorageProvider(clientId, clientSecret, region, bucket, endpoint)
	}
	return nil
}

var ossClient oss.StorageInterface

func GetOSSClient() oss.StorageInterface {
	return ossClient
}

func InitOSSClient(conf *conf.Conf) {
	ossClient = GetStorageProvider(conf.Oss.ProviderType, conf.Oss.AccessKeyId, conf.Oss.AccessKeySecret, "", conf.Oss.BucketName, conf.Oss.Endpoint)
}
