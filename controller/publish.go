package controller

import (
	"DouSheng/service"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

func Publish(c *gin.Context) {
	token := c.PostForm("token")
	// ckId就是当前用户的ID
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error()},
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title := c.Query("title")
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", ckId, filename)
	saveFile := filepath.Join("./public/", finalName)

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoUrl := fmt.Sprintf("http://192.168.1.100:8080/static/%s", finalName)
	if err := service.PublishVideo(&videoUrl, &title, ckId); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	token := c.Query("token")
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	userId_str := c.Query("user_id")
	userId, err := strconv.ParseInt(userId_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}

	videos, err := service.QueryUserPublishList(userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	// 接下来判断ckId用户是否喜爱这些视频
	for i := 0; i < len(videos); i++ {
		videos[i].IsFavorite = service.IsUserLikeVideo(ckId, videos[i].Id)
	}
	log.Printf("publish list size: %d", len(videos))
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
