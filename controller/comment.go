package controller

import (
	"DouSheng/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []service.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment service.Comment `json:"comment,omitempty"`
}

func CommentAction(c *gin.Context) {
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
	actionType := c.Query("action_type")

	if actionType == "1" {
		comment_text := c.Query("comment_text")
		comment, err := service.UserPublishComment(ckId, videoId, &comment_text)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
			},
			Comment: *comment,
		})
	} else if actionType == "2" {
		commentId_str := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentId_str, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Unknow comment ID",
			})
			return
		}
		err = service.UserDeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unknow action"})
	}
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	_, err := service.CheckTokenReturnID(&token)
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
			StatusMsg:  "Unknow video ID",
		})
		return
	}

	comments, err := service.QueryCommentListByVideoId(videoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
