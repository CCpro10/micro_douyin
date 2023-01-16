package domain

import (
	"context"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
)

type Comment struct {
	Id         int64   `json:"id"`
	User       *Author `json:"user"`
	Content    string  `json:"content"`
	CreateDate string  `json:"create_date"` //评论发布日期，格式 mm-dd
}

func FillComment(comment *repository.Comment, user *repository.User) *Comment {
	return &Comment{
		Id: comment.CommentId,
		User: &Author{
			Id:            user.UserId,
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true, //这里作者是自己,所以是true
		},
		Content:    comment.Content,
		CreateDate: util.Timestamp2Date(comment.CreateTime), //mm-dd
	}
}

func FillCommentList(ctx context.Context, comments []*repository.Comment, userId int64) ([]*Comment, error) {
	userIds := GetUserIdsFromCommentList(comments)

	authors, err := getAuthorsFromIds(ctx, userIds, userId)
	if err != nil {
		return nil, err
	}
	commentDOs := make([]*Comment, len(comments))
	for i := range comments {
		commentDOs[i] = &Comment{
			Id:         comments[i].CommentId,
			User:       authors[i],
			Content:    comments[i].Content,
			CreateDate: util.Timestamp2Date(comments[i].CreateTime), //mm-dd
		}
	}
	return commentDOs, nil
}

func GetUserIdsFromCommentList(comments []*repository.Comment) []int64 {
	ids := make([]int64, len(comments))
	for i := range comments {
		ids[i] = comments[i].UserId
	}
	return ids
}
