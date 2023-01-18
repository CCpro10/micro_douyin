package repository

import (
	"context"
	"gorm.io/gorm"
)

type Video struct {
	VideoId       int64  `json:"video_id" gorm:"primaryKey; bigint unsigned; not null"`
	UserId        int64  `json:"user_id" gorm:"bigint unsigned; not null"`
	PlayUrl       string `json:"play_url" gorm:"varchar(511); not null; default: '' "`
	CoverUrl      string `json:"cover_url" gorm:"varchar(511); not null; default: '' "`
	Title         string `json:"title" gorm:"varchar(255); not null"`
	FavoriteCount int64  `json:"favorite_count" gorm:"not null"`
	CommentCount  int64  `json:"comment_count" gorm:"not null"`
	CreateTime    int64  `json:"create_time" gorm:"autoCreateTime; not null"`
	ModifyTime    int64  `json:"modify_time" gorm:"autoUpdateTime; not null"`
	DeleteTime    int64  `json:"delete_time" gorm:"not null"`
}

type IVideoRepository interface {
	Create(context.Context, *Video) error
	FindByUserId(context.Context, int64) ([]*Video, error)
	FindByVideoId(context.Context, int64) (*Video, error)
	FindWithLimit(context.Context, int) ([]*Video, error)
	FindByCreateTimeWithLimit(context.Context, int64, int) ([]*Video, error)
	AddVideoFavoriteCount(context.Context, int64) error
	ReduceVideoFavoriteCount(context.Context, int64) error
	FindByVideoIds(context.Context, []int64) ([]*Video, error)
}
type VideoRepository struct{}

func (r *VideoRepository) Create(ctx context.Context, video *Video) error {
	return DB.WithContext(ctx).Create(&video).Error
}

func (r *VideoRepository) FindByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Where("user_id = ? and delete_time = 0", userId).Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) FindByVideoId(ctx context.Context, videoId int64) (*Video, error) {
	video := Video{}
	err := DB.WithContext(ctx).Where("video_id = ? and delete_time = 0", videoId).First(&video).Error
	return &video, err
}

func (r *VideoRepository) FindWithLimit(ctx context.Context, limit int) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("delete_time = 0").Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) FindByCreateTimeWithLimit(ctx context.Context, createTime int64, limit int) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("create_time < ? and delete_time = 0", createTime).Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) UpdadteFavoriteCount(ctx context.Context, video *Video, videoId int64) (err error) {
	err = DB.WithContext(ctx).Model(&Video{}).Where("video_id = ? ", videoId).Update("favorite_count", video.FavoriteCount).Error
	return err
}

func (r *VideoRepository) AddVideoFavoriteCount(ctx context.Context, videoId int64) error {
	return DB.WithContext(ctx).Model(&Video{}).Where("video_id = ? and delete_time = 0", videoId).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

func (r *VideoRepository) ReduceVideoFavoriteCount(ctx context.Context, videoId int64) error {
	return DB.WithContext(ctx).Model(&Video{}).Where("video_id = ?", videoId).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
}

func (r *VideoRepository) FindByVideoIds(ctx context.Context, videoList []int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Model(&Video{}).Where("video_id in (?) and delete_time = 0", videoList).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
