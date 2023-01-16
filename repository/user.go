package repository

import (
	"context"
	"time"
)

type User struct {
	UserId        int64  `json:"user_id" gorm:"primaryKey"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	CreateTime    int64  `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime    int64  `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime    int64  `json:"delete_time"`
}

type IUserRepository interface {
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	DeleteByUserId(context.Context, int64) error
	FindByUserId(context.Context, int64) (*User, error)
	FindByUserIds(context.Context, []int64) ([]*User, error)
	FindByUsername(context.Context, string) (*User, error)
}

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return DB.WithContext(ctx).Create(&user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return DB.WithContext(ctx).Where("delete_time = 0").Updates(&user).Error
}

func (r *UserRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", userId).Update("delete_time", time.Now().Unix()).Error
}

func (r *UserRepository) FindByUserId(ctx context.Context, userId int64) (user *User, err error) {
	err = DB.WithContext(ctx).Where("user_id = ? and delete_time = 0", userId).First(&user).Error
	return user, err
}

func (r *UserRepository) FindByUserIds(ctx context.Context, userIdList []int64) ([]*User, error) {
	users := make([]*User, 0)
	err := DB.WithContext(ctx).Where("user_id in (?) and delete_time = 0", userIdList).Find(&users).Error
	if err != nil {
		return nil, err
	}
	// record: id -> user
	id2UserMap := map[int64]*User{}
	for i := range users {
		if users[i] != nil {
			id2UserMap[users[i].UserId] = users[i]
		}
	}
	// fill
	res := make([]*User, len(userIdList))
	for i, id := range userIdList {
		res[i] = &User{
			UserId:        id2UserMap[id].UserId,
			Username:      id2UserMap[id].Username,
			Password:      id2UserMap[id].Password,
			FollowCount:   id2UserMap[id].FollowCount,
			FollowerCount: id2UserMap[id].FollowerCount,
			CreateTime:    id2UserMap[id].CreateTime,
			ModifyTime:    id2UserMap[id].ModifyTime,
			DeleteTime:    id2UserMap[id].DeleteTime,
		}
	}
	return res, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (user *User, err error) {
	err = DB.WithContext(ctx).Where("username = ? and delete_time = 0", username).First(&user).Error
	return user, err
}
