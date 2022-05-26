package service

import (
	"DouSheng/database"
	"errors"
)

// 获取Feed视频流
func GetFeed() ([]Video, error) {
	videosDB, err := database.NewDaoInstance().QueryVideos()
	if err != nil {
		return nil, err
	}
	videosLen := len(videosDB)
	videosCtr := make([]Video, videosLen)
	for i := 0; i < videosLen; i++ {
		favCount := database.NewDaoInstance().CountVideoLikesByVideoId(videosDB[i].VideoId)
		comCount := database.NewDaoInstance().CountVideoCommentsByVideoId(videosDB[i].VideoId)
		_, videoPublisher := QueryUserByUserId(videosDB[i].Publisher)
		// isLike := database.NewDaoInstance().QueryIsUserLikeVideo(userId, videosDB[i].VideoId)
		videosCtr[i] = Video{
			Id:            videosDB[i].VideoId,
			Author:        *videoPublisher,
			PlayUrl:       videosDB[i].PlayUrl,
			CoverUrl:      videosDB[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			IsFavorite:    false,
		}
	}
	return videosCtr, nil
}

// 发布视频
func PublishVideo(playUrl *string, user *User) error {
	videoDB := database.Video{
		Publisher: user.Id,
		PlayUrl:   *playUrl,
	}
	return database.NewDaoInstance().InsertVideo(&videoDB)
}

/* 查找用户发布了的视频
   select * from `Video` Where publisher = user_id */
func QueryUserPublishList(userId int64) ([]Video, error) {
	exist, user := QueryUserByUserId(userId)
	if !exist {
		return nil, errors.New("User doesn't exist")
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
		isLike := database.NewDaoInstance().QueryIsUserLikeVideo(userId, (*videosDB)[i].VideoId)
		videosCtr[i] = Video{
			Id:            (*videosDB)[i].VideoId,
			Author:        *user,
			PlayUrl:       (*videosDB)[i].PlayUrl,
			CoverUrl:      (*videosDB)[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			IsFavorite:    isLike,
		}
	}
	return videosCtr, nil
}

/* 查找用户点赞了的视频
   select * from `Video` Where publisher = user_id */
func QueryUserLikeList(userId int64) ([]Video, error) {
	// exist, user := QueryUserByUserId(userId)
	// if !exist {
	// 	return nil, errors.New("User doesn't exist")
	// }
	videosDB, err := database.NewDaoInstance().QueryVideosLikedByUserId(userId)
	if err != nil {
		return nil, err
	}
	videosLen := len(*videosDB)
	videosCtr := make([]Video, videosLen)
	for i := 0; i < videosLen; i++ {
		favCount := database.NewDaoInstance().CountVideoLikesByVideoId((*videosDB)[i].VideoId)
		comCount := database.NewDaoInstance().CountVideoCommentsByVideoId((*videosDB)[i].VideoId)
		_, videoPublisher := QueryUserByUserId((*videosDB)[i].Publisher)
		isLike := database.NewDaoInstance().QueryIsUserLikeVideo(userId, (*videosDB)[i].VideoId)
		videosCtr[i] = Video{
			Id:            (*videosDB)[i].VideoId,
			Author:        *videoPublisher,
			PlayUrl:       (*videosDB)[i].PlayUrl,
			CoverUrl:      (*videosDB)[i].CoverUrl,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			IsFavorite:    isLike,
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
