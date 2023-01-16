package service

import (
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func Register(c *gin.Context, username string, password string) (userId int64, token string, err error) {

	//判断用户名是否已被使用
	_, err = repository.GetUserRepository().FindByUsername(c, username)
	if err == nil { //用户名已被使用
		err = util.ErrUserExisted
		log.Println(err)
		return
	} else if err != gorm.ErrRecordNotFound { //出现了其他错误
		log.Println(err)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}

	//创建用户实例,存入注册信息
	user := repository.User{
		UserId:   util.GenerateId(),
		Username: username,
		Password: string(hashPassword),
	}
	err = repository.GetUserRepository().Create(c, &user)
	if err != nil {
		log.Println(err)
		return
	}

	//生成token
	if token, err = util.GenerateToken(user.UserId); err != nil {
		log.Println(err)
		return
	}
	return user.UserId, token, nil
}

func Login(c *gin.Context, username string, password string) (userId int64, token string, err error) {

	//通过用户名查找信息
	user, err := repository.GetUserRepository().FindByUsername(c, username)
	if err != nil {
		log.Println(err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err)
		err = util.ErrWrongPassword
		return
	}

	//生成token
	if token, err = util.GenerateToken(user.UserId); err != nil {
		log.Println(err)
		return
	}
	return user.UserId, token, nil
}
