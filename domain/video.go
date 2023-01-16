package domain

import (
	"context"
	"github.com/CCpro10/micro_douyin/repository"
)

type VideoDO struct {
	Id            int64   `json:"id"`
	Author        *Author `json:"author"`
	PlayUrl       string  `json:"play_url"`
	CoverUrl      string  `json:"cover_url"`
	Title         string  `json:"title"`
	FavoriteCount int64   `json:"favorite_count"`
	CommentCount  int64   `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
}

func FillVideoList(ctx context.Context, videoList []*repository.Video, audienceId int64, isFavorite bool) ([]*VideoDO, error) {
	res := make([]*VideoDO, len(videoList))

	authorIds := GetUserIdsFromVideoList(videoList)
	authors, err := getAuthorsFromIds(ctx, authorIds, audienceId)
	if err != nil {
		return nil, err
	}

	// TODO audienceId 用来查关注和点赞

	for i := range videoList {
		//isFavorite := flag
		res[i] = &VideoDO{
			Id:            videoList[i].VideoId,
			Author:        authors[i],
			PlayUrl:       videoList[i].PlayUrl,
			CoverUrl:      videoList[i].CoverUrl,
			Title:         videoList[i].Title,
			FavoriteCount: videoList[i].FavoriteCount, //TODO
			CommentCount:  videoList[i].CommentCount,  //TODO
			IsFavorite:    isFavorite,                 //TODO
		}
	}

	return res, nil
}

func GetUserIdsFromVideoList(videoList []*repository.Video) []int64 {
	ids := make([]int64, len(videoList))
	for i := range videoList {
		ids[i] = videoList[i].UserId
	}
	return ids
}

func getAuthorsFromIds(ctx context.Context, ids []int64, audienceId int64) ([]*Author, error) {
	users, err := repository.GetUserRepository().FindByUserIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	followedUserSet, err := getFollowedUserSet(ctx, audienceId)
	if err != nil {
		return nil, err
	}

	authors := make([]*Author, len(ids))
	for i := range users {
		_, ok := followedUserSet[users[i].UserId]
		authors[i] = FillAuthor(users[i], ok)
	}

	return authors, nil
}

// 获取该用户关注的用户的ID集合
func getFollowedUserSet(ctx context.Context, fromId int64) (map[int64]struct{}, error) {
	followedUsers, err := repository.GetFollowRepository().FindByFromUserId(ctx, fromId)
	if err != nil {
		return nil, err
	}
	followedUserSet := map[int64]struct{}{}
	for i := range followedUsers {
		followedUserSet[followedUsers[i].ToUserId] = struct{}{}
	}
	// 默认让前端显示关注自己
	followedUserSet[fromId] = struct{}{}
	return followedUserSet, nil
}
