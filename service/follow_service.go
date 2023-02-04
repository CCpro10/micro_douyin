package service

import (
	"context"
	"errors"
	"log"

	"github.com/CCpro10/micro_douyin/domain"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	FollowAction   = "1"         // 1-关注
	UnFollowAction = "2"         // 2-取消关注
	Following      = "following" // 关注
	Follower       = "follower"  // 粉丝
)

func RelationAction(c *gin.Context, userId int64, toUserId int64, actionType string) error {
	switch actionType {
	case FollowAction:
		return Follow(c, userId, toUserId)
	case UnFollowAction:
		return UnFollow(c, userId, toUserId)
	}
	return nil
}

// Follow 关注
func Follow(ctx context.Context, userId, toUserId int64) error {
	// error如果为空，说明已经关注了
	err := repository.GetFollowRepository().FindByUserId(ctx, userId, toUserId)
	if err == nil {
		return util.ErrIsFollow
	}
	// ErrRecordNotFound 表示查询不到记录的错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Create(ctx, userId, toUserId); err != nil {
		return err
	}
	// to_user的User对象“粉丝总数” follower_count + 1
	if err := repository.GetFollowRepository().AddFollowerCount(ctx, toUserId); err != nil {
		return err
	}
	// user的User对象“关注总数” follow_count + 1
	if err := repository.GetFollowRepository().AddFollowCount(ctx, userId); err != nil {
		return err
	}
	return nil
}

// UnFollow 取消关注
func UnFollow(ctx *gin.Context, userId, toUserId int64) error {
	// error为空表示已查询到对应记录，继续取消关注逻辑，其余错误均返回
	err := repository.GetFollowRepository().FindByUserId(ctx, userId, toUserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return util.ErrNotFollow
	}
	if err != nil {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Delete(ctx, userId, toUserId); err != nil {
		return err
	}
	// to_user的User对象“粉丝总数” follower_count + 1
	if err := repository.GetFollowRepository().ReduceFollowerCount(ctx, toUserId); err != nil {
		return err
	}
	// user的User对象“关注总数” follow_count + 1
	if err := repository.GetFollowRepository().ReduceFollowCount(ctx, userId); err != nil {
		return err
	}
	return nil
}

func IsFollowed(ctx *gin.Context, userId, toUserId int64) bool {
	err := repository.GetFollowRepository().FindByUserId(ctx, userId, toUserId)
	if err == nil {
		return false
	} else {
		return true
	}
}

// GetFollowList 获得关注列表
func GetFollowList(ctx *gin.Context, userId int64) []*domain.UserInfo {
	followList, err := repository.GetFollowRepository().FindByFromUserId(ctx, userId)
	if err != nil {
		log.Println("FindByFromUserId Failed", err)
		return nil
	}

	userIds := domain.GetToUserIds(followList)
	userList, err := GetUserInfosByIds(ctx, userIds, userId, Following)
	if err != nil {
		log.Println("GetUserInfosByIds Failed", err)
		return nil
	}
	return userList
}

// GetFollowerList 获得粉丝列表
func GetFollowerList(ctx *gin.Context, userId int64) []*domain.UserInfo {
	followList, err := repository.GetFollowRepository().FindByToUserId(ctx, userId)
	if err != nil {
		log.Println("FindByToUserId Failed", err)
		return nil
	}

	userIds := domain.GetFromUserIds(followList)
	userList, err := GetUserInfosByIds(ctx, userIds, userId, Follower)
	if err != nil {
		log.Println("GetUserInfosByIds Failed", err)
		return nil
	}

	return userList
}
