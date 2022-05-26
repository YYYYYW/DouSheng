package controller

import (
	"DouSheng/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	// TODO finish feed
	videos, err := service.GetFeed()
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			NextTime: time.Now().Unix(),
		})
		return
	}
	token := c.Query("token")
	if token != "" {
		ckId, err := service.CheckTokenReturnID(&token)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		}
		for i := 0; i < len(videos); i++ {
			videos[i].Author.IsFollow = service.IsUserFollowToUser(ckId, videos[i].Author.Id)
			videos[i].IsFavorite = service.IsUserLikeVideo(ckId, videos[i].Id)
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
