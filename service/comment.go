package service

import (
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
)

func PublishComment(c *gin.Context, userId int64, videoId int64, content string) (*repository.Comment, *repository.User, error) {
	//检查是否含敏感词
	find, _ := util.Filter.FindIn(content)
	if find {
		return nil, nil, util.ErrSensitiveComment
	}
	comment := &repository.Comment{
		CommentId: util.GenerateId(),
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
	}
	err := repository.GetCommentRepository().Create(c, comment)
	if err != nil {
		return nil, nil, err
	}

	err = repository.GetCommentRepository().AddCommentCount(c, comment)
	if err != nil {
		return nil, nil, err
	}
	//此处userId是肯定存在的
	user, err := repository.GetUserRepository().FindByUserId(c, userId)
	if err != nil {
		return nil, nil, err
	}

	return comment, user, nil
}

func DeleteComment(c *gin.Context, commentId int64) (err error) {
	comment, err := repository.GetCommentRepository().FindByCommentId(c, commentId)
	if err != nil {
		return err
	}

	//先把video对应的comment_count减一,再删除评论
	err = repository.GetCommentRepository().ReduceCommentCount(c, comment)
	if err != nil {
		return err
	}

	err = repository.GetCommentRepository().DeleteByCommentId(c, commentId)
	if err != nil {
		return err
	}

	return nil
}
