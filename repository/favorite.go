package repository

import (
	"context"
	"time"
)

type Favorite struct {
	FavoriteId int64 `json:"favorite_id" gorm:"primaryKey"`
	UserId     int64 `json:"user_id" gorm:"index"`
	VideoId    int64 `json:"video_id" gorm:"index"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime;index"`
	ModifyTime int64 `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64 `json:"delete_time"`
}

type IFavoriteRepository interface {
	Create(context.Context, *Favorite) error
	Update(context.Context, *Favorite) error
	DeleteByUserIdAndVideoId(context.Context, int64, int64) error
	FindVideoListByUserId(context.Context, int64) ([]*Favorite, error)
	FindFavoriteRecord(ctx context.Context, userId int64, videoId int64) (favorite *Favorite, err error)
	DeleteFavorite(ctx context.Context, userId int64, videoId int64, favorite *Favorite) (err error)
}

type FavoriteRepository struct{}

func (r *FavoriteRepository) FindVideoListByUserId(ctx context.Context, userId int64) ([]*Favorite, error) {
	favorite := make([]*Favorite, 0)
	err := DB.WithContext(ctx).Where("user_id = ? and delete_time = 0", userId).Find(&favorite).Error
	return favorite, err
}

func (r *FavoriteRepository) FindFavoriteRecord(ctx context.Context, userId int64, videoId int64) (favorite *Favorite, err error) {
	err = DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).First(&favorite).Error
	return favorite, err
}

func (r *FavoriteRepository) DeleteFavorite(ctx context.Context, userId int64, videoId int64, favorite *Favorite) (err error) {
	err = DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Update("delete_time", favorite.DeleteTime).Error
	return err
}

func (r *FavoriteRepository) Create(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Create(&favorite).Error
}

func (r *FavoriteRepository) Update(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Where("user_id = ? and ").Updates(&favorite).Error
}

func (r *FavoriteRepository) DeleteByUserIdAndVideoId(ctx context.Context, userId int64, videoId int64) error {
	return DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Update("delete_time", time.Now().Unix()).Error
}
