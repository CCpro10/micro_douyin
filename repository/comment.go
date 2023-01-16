package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	CommentId  int64  `json:"comment_id"`
	UserId     int64  `json:"user_id"`
	VideoId    int64  `json:"video_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime int64  `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64  `json:"delete_time"`
}

type ICommentRepository interface {
	Create(context.Context, *Comment) error
	FindByCommentId(context.Context, int64) (*Comment, error)
	FindByVideoId(context.Context, int64) ([]*Comment, error)
	AddCommentCount(context.Context, *Comment) error
	ReduceCommentCount(context.Context, *Comment) error
	DeleteByCommentId(context.Context, int64) error
}

type CommentRepository struct{}

func (r *CommentRepository) Create(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Create(&comment).Error
}

func (r *CommentRepository) FindByCommentId(ctx context.Context, commentId int64) (*Comment, error) {
	comment := Comment{}
	err := DB.WithContext(ctx).Where("comment_id = ? and delete_time = 0", commentId).First(&comment).Error
	return &comment, err
}

func (r *CommentRepository) AddCommentCount(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Model(&Video{}).Where("video_id = ? and delete_time = 0", comment.VideoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (r *CommentRepository) ReduceCommentCount(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Model(&Video{}).Where("video_id = ? and delete_time = 0", comment.VideoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error
}

func (r *CommentRepository) DeleteByCommentId(ctx context.Context, commentId int64) error {
	err := DB.Model(&Comment{}).WithContext(ctx).Where("comment_id =? and delete_time =0", commentId).Update("delete_time", time.Now().Unix()).Error
	return err
}

func (r *CommentRepository) FindByVideoId(ctx context.Context, videoId int64) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Where("video_id = ? and delete_time =0", videoId).Find(&comments).Error
	return comments, err
}
