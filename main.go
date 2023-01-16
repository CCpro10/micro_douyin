package main

import (
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/CCpro10/micro_douyin/handler"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
)

const (
	ConfPath = "./conf/conf.yaml"
)

func main() {
	conf.Init(ConfPath)
	repository.InitDB()
	util.InitIdGenerator()
	util.InitJWTVal()
	util.InitValidate()
	util.InitOSSClient(conf.Config)
	util.InitFilter()

	r := gin.Default()
	handler.Register(r)
	r.Run(conf.Config.Server.Port)
}
