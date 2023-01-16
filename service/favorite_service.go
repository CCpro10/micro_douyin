package service

import (
	"context"
	"errors"
	"github.com/CCpro10/micro_douyin/domain"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

const (
	likeAction   = "1"
	unLikeAction = "2"
)

func FavoriteAction(c *gin.Context, userId int64, videoId int64, actionType string) error {
	switch actionType {
	case likeAction:
		return Like(c, userId, videoId)
	case unLikeAction:
		return Unlike(c, userId, videoId)
	}
	return util.ErrParamError
}

//Like
func Like(ctx context.Context, userId int64, videoId int64) error {
	//find favorite record by userId and videoId
	favorite, err := repository.GetFavoriteRepository().FindFavoriteRecord(ctx, userId, videoId)
	//error is not "not found" ,return err
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	//record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		favoriteId := util.GenerateId()
		favorite := &repository.Favorite{
			FavoriteId: favoriteId,
			UserId:     userId,
			VideoId:    videoId,
		}
		//create favorite record
		if err = repository.GetFavoriteRepository().Create(ctx, favorite); err != nil {
			return err
		}
	} else {
		//update favorite record: delete_time = 0
		favorite.DeleteTime = 0
		if err = repository.GetFavoriteRepository().DeleteFavorite(ctx, userId, videoId, favorite); err != nil {
			return err
		}
	}
	//add video’s favoriteCount
	if err = repository.GetVideoRepository().AddVideoFavoriteCount(ctx, videoId); err != nil {
		return err
	}
	return nil
}

//Unlike
func Unlike(ctx context.Context, userId int64, videoId int64) error {
	// update delete_time = now()
	if err := repository.GetFavoriteRepository().DeleteByUserIdAndVideoId(ctx, userId, videoId); err != nil {
		return err
	}
	//decrease video's favoriteCount
	if err := repository.GetVideoRepository().ReduceVideoFavoriteCount(ctx, videoId); err != nil {
		return err
	}
	return nil
}

//get favoriteList
func GetFavoriteList(ctx context.Context, userId int64) ([]*domain.VideoDO, error) {
	var (
		videoIdList []int64
	)

	favoriteList, err := repository.GetFavoriteRepository().FindVideoListByUserId(ctx, userId)
	if err != nil {
		log.Println("GetVideoRepository().FindVideoListByUserId Failed", err)
		return nil, err
	}

	for _, x := range favoriteList {
		videoIdList = append(videoIdList, x.VideoId)
	}
	videoList, err := repository.GetVideoRepository().FindByVideoIds(ctx, videoIdList)
	if err != nil {
		log.Println("GetFavoriteList FindByVideoIds Failed", err)
		return nil, err
	}

	videoDOs, err := domain.FillVideoList(ctx, videoList, userId, true)
	if err != nil {
		log.Println("GetFavoriteList FillVideoList Failed", err)
		return nil, err
	}

	return videoDOs, nil
}
