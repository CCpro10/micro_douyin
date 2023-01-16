package repository

import (
	"context"
	"errors"
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"testing"
)

const (
	confPath = "../conf/conf.yaml"
)

func TestUserRepo(t *testing.T) {
	conf.Init(confPath)
	initMySQL()
	initRepository()
	util.InitIdGenerator()

	ctx := context.Background()
	userId := util.GenerateId()
	newUser := &User{
		UserId:   userId,
		Username: "qizong007",
		Password: "123456",
	}

	if err := GetUserRepository().Create(ctx, newUser); err != nil {
		log.Println(err)
		return
	}
	log.Println(newUser)

	user, err := GetUserRepository().FindByUserId(ctx, userId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user)

	user.Password = "6666"
	if err = GetUserRepository().Update(ctx, user); err != nil {
		log.Println(err)
		return
	}
	log.Println(user)

	if err = GetUserRepository().DeleteByUserId(ctx, userId); err != nil {
		log.Println(err)
		return
	}

	_, err = GetUserRepository().FindByUserId(ctx, userId)
	assert.Equal(t, true, errors.Is(err, gorm.ErrRecordNotFound))

	id1 := util.GenerateId()
	id2 := util.GenerateId()
	id3 := util.GenerateId()

	user1 := &User{
		UserId:   id1,
		Username: "user1",
		Password: "123456",
	}

	user2 := &User{
		UserId:   id2,
		Username: "user2",
		Password: "123456",
	}

	user3 := &User{
		UserId:   id3,
		Username: "user3",
		Password: "123456",
	}

	_ = GetUserRepository().Create(ctx, user1)
	_ = GetUserRepository().Create(ctx, user2)
	_ = GetUserRepository().Create(ctx, user3)

	users, err := GetUserRepository().FindByUserIds(ctx, []int64{id1, id2, id3})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(users)

	user, err = GetUserRepository().FindByUsername(ctx, "user2")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user)
}
