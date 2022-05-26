package controller

import (
	"DouSheng/service"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	exist, id := service.QueryUserIdByToken(&token)
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	exist, user := service.QueryUserByUserId(id)
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
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

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoUrl := fmt.Sprintf("http://192.168.1.100:8080/static/%s", finalName)
	if err := service.PublishVideo(&videoUrl, user); err != nil {
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

	// if _, exist := usersLoginInfo[token]; !exist {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }

	// data, err := c.FormFile("data")
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  err.Error(),
	// 	})
	// 	return
	// }

	// filename := filepath.Base(data.Filename)
	// user := usersLoginInfo[token]
	// finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// saveFile := filepath.Join("./public/", finalName)
	// if err := c.SaveUploadedFile(data, saveFile); err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  err.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, Response{
	// 	StatusCode: 0,
	// 	StatusMsg:  finalName + " uploaded successfully",
	// })
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	id_str := c.Query("user_id")
	id, err := strconv.ParseInt(id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
	}

	token := c.Query("token")
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
	videos, err := service.QueryUserPublishList(id)
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
