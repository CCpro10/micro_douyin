package domain

import (
	"github.com/CCpro10/micro_douyin/repository"
)

// Author 作为观众看到的创作者
type Author struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func FillAuthor(user *repository.User, isFollow bool) *Author {
	author := &Author{
		Id:            user.UserId,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	return author
}
