package middleware

import (
	"errors"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"log"
)

func JWT(c *gin.Context) {

	token := c.Query("token")
	//解析token
	userId, err := util.ParseToken(token)
	if err != nil { //ParseToken只会返回两种错误
		if errors.Is(err, util.ErrNoAuth) {
			log.Println("JWTMiddleWare Token <Nil>")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.NoAuth,
			})
			c.Abort()
			return
		}
		if errors.Is(err, util.ErrWrongAuth) {
			log.Println("JWTMiddleWare Token Wrong ,Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
			c.Abort()
			return
		}
	}
	c.Set("userId", userId)
	c.Next()
}

//获取JWT middleware解析过出的userId,
//未经过JWT middleware解析时被调用,会返回错误
func GetUserId(ctx *gin.Context) (int64, error) {
	userId, exists := ctx.Get("userId")
	if exists {
		return userId.(int64), nil
	}
	return 0, errors.New("UserId Not Set In Middleware")
}
