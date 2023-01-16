package handler

import (
	"errors"
	"github.com/CCpro10/micro_douyin/domain"
	"github.com/CCpro10/micro_douyin/middleware"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/service"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"strings"
)

type RegisterReq struct {
	Username string `validate:"required,max=32"`
	Password string `validate:"required,min=6,max=32"`
}

func RegisterHandler(c *gin.Context) {
	//解析,处理参数
	var req RegisterReq
	var resp util.HttpResponse
	req.Username = c.Query("username")
	req.Password = c.Query("password")
	req.Username = strings.TrimSpace(req.Username) //去首尾空格

	err := util.Validate.Struct(req) // 执行验证
	if err != nil {
		resp.StatusCode = util.ParamError
		util.MakeResponse(c, &resp)
		return
	}
	userId, token, err := service.Register(c, req.Username, req.Password)
	resp.ReturnVal = map[string]interface{}{
		"user_id": userId,
		"token":   token,
	}
	if err != nil {
		if errors.Is(err, util.ErrUserExisted) { //用户已存在
			resp.StatusCode = util.UserExisted
			util.MakeResponse(c, &resp)
			return
		}

		resp.StatusCode = util.InternalServerError
		util.MakeResponse(c, &resp)
		return
	}
	resp.StatusCode = util.Success
	util.MakeResponse(c, &resp)
	return
}

func LoginHandler(c *gin.Context) {
	//登录请求和注册请求是一样的
	var req RegisterReq
	var resp util.HttpResponse
	req.Username = c.Query("username")
	req.Password = c.Query("password")
	//去username首尾空格
	req.Username = strings.TrimSpace(req.Username)

	validate := validator.New() // 创建验证器
	err := validate.Struct(req) // 执行验证
	if err != nil {
		resp.StatusCode = util.ParamError
		util.MakeResponse(c, &resp)
		return
	}

	userId, token, err := service.Login(c, req.Username, req.Password)
	resp.ReturnVal = map[string]interface{}{
		"user_id": userId,
		"token":   token,
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //用户未注册
			resp.StatusCode = util.UserNotExist
			util.MakeResponse(c, &resp)
			return
		}

		if errors.Is(err, util.ErrWrongPassword) { //用户密码错误
			resp.StatusCode = util.WrongPassword
			util.MakeResponse(c, &resp)
			return
		}

		resp.StatusCode = util.InternalServerError
		util.MakeResponse(c, &resp)
		return
	}
	resp.StatusCode = util.Success
	util.MakeResponse(c, &resp)
	return
}

func GetUserInfoHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	loginUserId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	userIdStr := c.Query("user_id")
	userId, err := util.Str2Int64(userIdStr)
	if err != nil {
		log.Println("Str2Int64 Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	user, err := repository.GetUserRepository().FindByUserId(c, userId)
	if err != nil {
		log.Println("GetUserRepository().FindByUserId Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	isFollow := true

	if loginUserId == userId {
		err = repository.GetFollowRepository().FindByUserId(c, loginUserId, userId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("GetFollowRepository().FindByUserId Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
		if err != nil {
			isFollow = false
		}
	}

	author := domain.FillAuthor(user, isFollow)
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"user": author,
		},
	})
}
