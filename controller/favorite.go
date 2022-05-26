package controller

import (
	"DouSheng/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 前端传token，video_id，action_type。没有userId
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	videoId_str := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoId_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}
	action_type := c.Query("action_type")

	if action_type == "1" {
		err := service.UserLikeVideo(ckId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" {
		err := service.UserUnLikeVideo(ckId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unknow action"})
	}
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	userId_str := c.Query("user_id")
	userId, err := strconv.ParseInt(userId_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow user ID",
		})
	}

	videos, err := service.QueryUserLikeList(userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// 判断ckId用户是否喜欢这些视频
	if ckId == userId {
		for i := 0; i < len(videos); i++ {
			videos[i].IsFavorite = true
		}
	} else {
		for i := 0; i < len(videos); i++ {
			videos[i].IsFavorite = service.IsUserLikeVideo(ckId, videos[i].Id)
		}
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
