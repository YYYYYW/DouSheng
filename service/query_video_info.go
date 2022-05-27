package service

import (
	"DouSheng/database"
	"log"
	"time"
)

// 获取Feed视频流
func GetFeed(latestTime int64) ([]Video, int64, error) {
	log.Printf("time now: %d, latestTime: %d", time.Now().UnixMilli(), latestTime)
	videosDB, err := database.NewDaoInstance().QueryVideos(latestTime)
	if err != nil {
		return nil, latestTime, err
	}
	videosLen := len(*videosDB)
	videosCtr := make([]Video, videosLen)
	for i := 0; i < videosLen; i++ {
		favCount := database.NewDaoInstance().CountVideoLikesByVideoId((*videosDB)[i].VideoId)
		comCount := database.NewDaoInstance().CountVideoCommentsByVideoId((*videosDB)[i].VideoId)
		videoPublisher, _ := QueryUserByUserId((*videosDB)[i].Publisher)
		videosCtr[i] = Video{
			Id:            (*videosDB)[i].VideoId,
			Author:        *videoPublisher,
			PlayUrl:       (*videosDB)[i].PlayUrl,
			CoverUrl:      (*videosDB)[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			IsFavorite:    false,
		}
	}
	nextTime := latestTime
	if videosLen > 0 {
		nextTime = (*videosDB)[0].CreateTime
	}
	return videosCtr, nextTime, nil
}

// 发布视频
func PublishVideo(playUrl *string, picUrl *string, title *string, userId int64) error {
	videoDB := database.Video{
		Publisher:  userId,
		PlayUrl:    *playUrl,
		CoverUrl:   *picUrl,
		Title:      *title,
		CreateTime: time.Now().UnixMilli(),
	}
	return database.NewDaoInstance().InsertVideo(&videoDB)
}

/* 查找用户发布了的视频
   select * from `Video` Where publisher = user_id */
func QueryUserPublishList(userId int64) ([]Video, error) {
	user, err := QueryUserByUserId(userId)
	if err != nil {
		return nil, err
	}

	videosDB, err := database.NewDaoInstance().QueryVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	videosLen := len(*videosDB)
	videosCtr := make([]Video, videosLen)
	for i := 0; i < videosLen; i++ {
		favCount := database.NewDaoInstance().CountVideoLikesByVideoId((*videosDB)[i].VideoId)
		comCount := database.NewDaoInstance().CountVideoCommentsByVideoId((*videosDB)[i].VideoId)
		// 查看是否喜欢应该是当前用户是否喜欢，而不是发布用户是否喜欢
		// isLike := database.NewDaoInstance().QueryIsUserLikeVideo(userId, (*videosDB)[i].VideoId)
		videosCtr[i] = Video{
			Id:            (*videosDB)[i].VideoId,
			Author:        *user,
			PlayUrl:       (*videosDB)[i].PlayUrl,
			CoverUrl:      (*videosDB)[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			Title:         (*videosDB)[i].Title,
		}
	}
	return videosCtr, nil
}

/* 查找用户点赞了的视频
   select * from `Video` Where publisher = user_id */
func QueryUserLikeList(userId int64) ([]Video, error) {
	videosDB, err := database.NewDaoInstance().QueryVideosLikedByUserId(userId)
	if err != nil {
		return nil, err
	}
	videosLen := len(*videosDB)
	videosCtr := make([]Video, videosLen)
	for i := 0; i < videosLen; i++ {
		favCount := database.NewDaoInstance().CountVideoLikesByVideoId((*videosDB)[i].VideoId)
		comCount := database.NewDaoInstance().CountVideoCommentsByVideoId((*videosDB)[i].VideoId)
		// TODO 如果发布者没有找到怎么办
		videoPublisher, _ := QueryUserByUserId((*videosDB)[i].Publisher)
		videosCtr[i] = Video{
			Id:            (*videosDB)[i].VideoId,
			Author:        *videoPublisher,
			PlayUrl:       (*videosDB)[i].PlayUrl,
			CoverUrl:      (*videosDB)[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			Title:         (*videosDB)[i].Title,
		}
	}
	return videosCtr, nil
}

// 用户点赞视频
func UserLikeVideo(userId int64, videoId int64) error {
	return database.NewDaoInstance().InsertLike(userId, videoId)
}

// 用户取消点赞视频
func UserUnLikeVideo(userId int64, videoId int64) error {
	return database.NewDaoInstance().DeleteLike(userId, videoId)
}

// 判断用户是否点赞了视频
func IsUserLikeVideo(userId int64, videoId int64) bool {
	return database.NewDaoInstance().QueryIsUserLikeVideo(userId, videoId)
}
