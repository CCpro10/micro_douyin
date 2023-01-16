package handler

import (
	"errors"
	"fmt"
	"github.com/CCpro10/micro_douyin/domain"
	"github.com/CCpro10/micro_douyin/middleware"
	"github.com/CCpro10/micro_douyin/repository"
	"github.com/CCpro10/micro_douyin/service"
	"github.com/CCpro10/micro_douyin/util"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	FeedLimit = 30
)

func VideoPublishHandler(c *gin.Context) {
	token := c.PostForm("token")

	userId, err := util.ParseToken(token)
	if err != nil { //ParseToken只会返回两种错误
		if errors.Is(err, util.ErrNoAuth) {
			log.Println("VideoPublishHandler Token <Nil>")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.NoAuth,
			})
			return
		}
		if errors.Is(err, util.ErrWrongAuth) {
			log.Println("VideoPublishHandler Token Wrong ,Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
			return
		}
	}

	title := c.PostForm("title")
	if title == "" {
		log.Println("VideoPublishHandler Title <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	// 读取视频文件数据
	videoFileHeader, err := c.FormFile("data")
	if err != nil {
		log.Println("VideoPublishHandler FormFile Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	videoFile, err := videoFileHeader.Open()
	if err != nil {
		log.Println("VideoPublishHandler Open File Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	objectName := fmt.Sprintf("%d/%d_%s", userId, time.Now().Unix(), videoFileHeader.Filename)

	// 上传视频文件
	err = service.VideoPublish(objectName, videoFile)
	if err != nil {
		log.Println("VideoPublish Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	playUrl := service.VideoUploadUrlPrefix + objectName
	coverUrl := playUrl + service.VideoCoverSuffix

	if err = repository.GetVideoRepository().Create(c, &repository.Video{
		VideoId:  util.GenerateId(),
		UserId:   userId,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}); err != nil {
		log.Println("GetVideoRepository().Create Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})
}

func VideoPublishedListHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	loginUserId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	userIdStr := c.Query("user_id")
	userId, err := util.Str2Int64(userIdStr)
	if err != nil {
		log.Println("Str2Int64 Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	videoList, err := repository.GetVideoRepository().FindByUserId(c, userId)
	if err != nil {
		log.Println("VideoPublishedListHandler FindByUserId Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	videoDOs, err := domain.FillVideoList(c, videoList, loginUserId, false)

	if err != nil {
		log.Println("FillVideoList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
		},
	})
}

func VideoFeedHandler(c *gin.Context) {
	var (
		userId    int64
		err       error
		videoList []*repository.Video
	)

	latestTimeStr := c.Query("latest_time")
	token := c.Query("token")

	//解析token,这里允许token为空的未登录用户刷视频,忽略了NoAuth错误
	userId, err = util.ParseToken(token)
	if err != nil {
		if errors.Is(err, util.ErrWrongAuth) { //token不为空,但是解析出错,这里会提醒用户重新登陆
			log.Println("JWTMiddleWare Token Wrong ,Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
		}
	}

	if latestTimeStr == "" { // 没有传入 latest_time
		videoList, err = repository.GetVideoRepository().FindWithLimit(c, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	} else { // 传入了 latest_time
		latestTime, err := util.Str2Int64(latestTimeStr)
		if err != nil {
			log.Println("Str2Int64 Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}

		videoList, err = repository.GetVideoRepository().FindByCreateTimeWithLimit(c, latestTime, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindByCreateTimeWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	}

	videoDOs, err := domain.FillVideoList(c, videoList, userId, false)
	if err != nil {
		log.Println("FillVideoList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
			"next_time":  getMostEarlyTime(videoList),
		},
	})

}

// 获取本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
func getMostEarlyTime(videos []*repository.Video) int64 {
	if len(videos) == 0 {
		return time.Now().Unix()
	}
	return videos[len(videos)-1].CreateTime
}
