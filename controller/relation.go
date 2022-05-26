package controller

import (
	"DouSheng/service"
	"log"
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
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
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

	if action_type == "1" {
		err := service.UserRelationToUser(ckId, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
	} else if action_type == "2" {
		err := service.UserUnRelationToUser(ckId, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unknow action"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

func FollowList(c *gin.Context) {
	log.Printf("request to follow list")
	token := c.Query("token")
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	user_id_str := c.Query("user_id")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
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

	// 对于结果用户列表，依次显示ckId用户是否关注了
	if ckId == user_id {
		for i := 0; i < len(followList); i++ {
			followList[i].IsFollow = true
		}
	} else {
		for i := 0; i < len(followList); i++ {
			followList[i].IsFollow = service.IsUserFollowToUser(ckId, followList[i].Id)
		}
	}

	// 暂时不设置用户列表中每一个用户的关注数和被关注数
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}

func FollowerList(c *gin.Context) {
	log.Printf("request to follower list")
	token := c.Query("token")
	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	user_id_str := c.Query("user_id")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Unknow ID",
		})
		return
	}

	followerList, err := service.QueryUserFollowerList(user_id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}

	// 对于结果用户列表，依次显示ckId用户是否关注了
	for i := 0; i < len(followerList); i++ {
		followerList[i].IsFollow = service.IsUserFollowToUser(ckId, followerList[i].Id)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followerList,
	})
}
