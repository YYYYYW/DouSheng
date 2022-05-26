package controller

import (
	"DouSheng/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 前端传token，video_id，action_type。没有userId
func FavoriteAction(c *gin.Context) {
	// userId_str := c.Query("user_id")
	// userId, err := strconv.ParseInt(userId_str, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  "Unknow ID",
	// 	})
	// 	return
	// }
	token := c.Query("token")
	exist, userId := service.QueryUserIdByToken(&token)
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
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
		err := service.UserLikeVideo(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" {
		err := service.UserUnLikeVideo(userId, videoId)
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
	id_str := c.Query("user_id")
	token := c.Query("token")
	id, err := strconv.ParseInt(id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow user ID",
		})
	}
	log.Printf("FavoriteList id: %d, token: %s", id, token)
	ret := service.CheckToken(id, &token)
	if ret == -1 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	} else if ret == -2 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "authentication failed",
		})
		return
	}

	videos, err := service.QueryUserLikeList(id)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
