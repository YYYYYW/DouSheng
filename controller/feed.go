package controller

import (
	"DouSheng/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	latestTime_str := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(latestTime_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			NextTime: 0,
		})
		return
	}
	videos, nextTime, err := service.GetFeed(latestTime)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			NextTime: 0,
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
		NextTime:  nextTime,
	})
}
