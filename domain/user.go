package domain

import (
	"github.com/CCpro10/micro_douyin/repository"
)

// UserInfo 返回用户信息
type UserInfo struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Avatar        string `json:"avatar"`
}

// FillUserInfo 将user对象转换为UserInfo对象
func FillUserInfo(user *repository.User, isFollow bool) *UserInfo {
	userInfo := &UserInfo{
		Id:            user.UserId,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
		Avatar:        user.Avatar,
	}
	return userInfo
}

// GetToUserIds 从返回的follow对象中，获得关注id列表
func GetToUserIds(followList []*repository.Follow) []int64 {
	ids := make([]int64, len(followList))
	for i := range followList {
		ids[i] = followList[i].ToUserId
	}
	return ids
}

// GetFromUserIds 从返回的follow对象中，获得粉丝id列表
func GetFromUserIds(followList []*repository.Follow) []int64 {
	ids := make([]int64, len(followList))
	for i := range followList {
		ids[i] = followList[i].FromUserId
	}
	return ids
}
