package controller

import (
	"DouSheng/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []service.User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	user_id_str := c.Query("user_id")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}
	to_user_id_str := c.Query("to_user_id")
	to_user_id, err := strconv.ParseInt(to_user_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}
	action_type := c.Query("action_type")
	ret := service.CheckToken(user_id, &token)
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
	if action_type == "1" {
		err := service.UserRelationToUser(user_id, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" {
		err := service.UserUnRelationToUser(user_id, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id_str := c.Query("user_id")
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
	followList, err := service.QueryUserFollowList(user_id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}

func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user_id_str := c.Query("user_id")
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
	followList, err := service.QueryUserFollowerList(user_id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}
