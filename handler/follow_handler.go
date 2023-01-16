package handler

import (
	"github.com/CCpro10/micro_douyin/middleware"
	"github.com/CCpro10/micro_douyin/service"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"log"
)

type RelationReq struct {
	UserID     int64
	ToUserID   int64
	ActionType string
}

func FollowActionHandler(c *gin.Context) {
	var req RelationReq
	req.ToUserID, _ = util.Str2Int64(c.Query("to_user_id"))
	req.ActionType = c.Query("action_type")

	//获取从JWTMiddleware解析好的userId
	userId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}
	req.UserID = userId // 客户端的请求里没有接口文档里说的userId，直接通过token解析
	err = service.RelationAction(c, req.UserID, req.ToUserID, req.ActionType)
	if err != nil {
		log.Println("RelationAction() Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})

}

func GetFollowListHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	userId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"user_list": service.GetFollowList(c, userId),
		},
	})
}

func GetFollowerListHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	userId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"user_list": service.GetFollowerList(c, userId),
		},
	})
}
