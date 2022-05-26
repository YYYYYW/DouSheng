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

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	user_id_str := c.Query("user_id")
	token := c.Query("token")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}
	ret := service.CheckToken(user_id, &token)
	if ret == -1 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
		return
	} else if ret == -2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "authentication failed",
			},
		})
		return
	}

	actionType := c.Query("action_type")
	video_id := c.GetInt64("video_id")

	if actionType == "1" {
		comment_text := c.Query("comment_text")
		comment, err := service.UserPublishComment(user_id, video_id, &comment_text)
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
		comment_id := c.GetInt64("comment_id")
		err := service.UserDeleteComment(comment_id)
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

	// if user, exist := usersLoginInfo[token]; exist {
	// 	if actionType == "1" {
	// 		text := c.Query("comment_text")
	// 		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
	// 			Comment: Comment{
	// 				Id:         1,
	// 				User:       user,
	// 				Content:    text,
	// 				CreateDate: "05-01",
	// 			}})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// token := c.Query("token")
	video_id_str := c.Query("video_id")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow video ID",
		})
		return
	}

	comments, err := service.QueryCommentListByVideoId(video_id)
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
