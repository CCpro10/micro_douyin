package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Follow struct {
	FollowId   int64 `json:"follow_id" gorm:"primaryKey"`
	FromUserId int64 `json:"from_user_id"`
	ToUserId   int64 `json:"to_user_id"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime int64 `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64 `json:"delete_time"`
}

type IFollowRepository interface {
	FindByUserId(context.Context, int64, int64) error
	Create(context.Context, int64, int64) error
	AddFollowerCount(context.Context, int64) error
	AddFollowCount(context.Context, int64) error
	Delete(context.Context, int64, int64) error
	ReduceFollowerCount(context.Context, int64) error
	ReduceFollowCount(context.Context, int64) error
	FindByFromUserId(context.Context, int64) ([]*Follow, error)
	FindByToUserId(context.Context, int64) ([]*Follow, error)
}

type FollowRepository struct{}

func (r *FollowRepository) FindByUserId(ctx context.Context, fromUserId, toUserId int64) error {
	return DB.WithContext(ctx).First(&Follow{}, "to_user_id = ? and from_user_id = ? and delete_time = 0", toUserId, fromUserId).Error
}

func (r *FollowRepository) Create(ctx context.Context, fromUserId, toUserId int64) error {
	return DB.WithContext(ctx).Create(&Follow{ToUserId: toUserId, FromUserId: fromUserId}).Error
}

func (r *FollowRepository) AddFollowerCount(ctx context.Context, toUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
}

func (r *FollowRepository) AddFollowCount(ctx context.Context, fromUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", fromUserId).
		UpdateColumn("follow_count", gorm.Expr("follow_count + 1")).Error
}

func (r *FollowRepository) Delete(ctx context.Context, fromUserId, toUserId int64) error {
	return DB.WithContext(ctx).Model(&Follow{}).Where("to_user_id=? AND from_user_id=?", toUserId, fromUserId).
		Update("delete_time", time.Now().Unix()).Error
}

func (r *FollowRepository) ReduceFollowerCount(ctx context.Context, toUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
}

func (r *FollowRepository) ReduceFollowCount(ctx context.Context, fromUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", fromUserId).
		UpdateColumn("follow_count", gorm.Expr("follow_count - 1")).Error
}

func (r *FollowRepository) FindByFromUserId(ctx context.Context, fromUserId int64) ([]*Follow, error) {
	followList := make([]*Follow, 0)
	err := DB.WithContext(ctx).Where("from_user_id = ? and delete_time = 0", fromUserId).Find(&followList).Error
	return followList, err
}

func (r *FollowRepository) FindByToUserId(ctx context.Context, toUserId int64) ([]*Follow, error) {
	followList := make([]*Follow, 0)
	err := DB.WithContext(ctx).Where("to_user_id = ? and delete_time = 0", toUserId).Find(&followList).Error
	return followList, err
}
