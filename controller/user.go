package controller

import (
	"DouSheng/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
// var usersLoginInfo = map[string]User{
// 	"zhangleidouyin": {
// 		Id:            1,
// 		Name:          "zhanglei",
// 		FollowCount:   10,
// 		FollowerCount: 5,
// 		IsFollow:      true,
// 	},
// }

// var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if _, err := service.QueryUserExisted(&username, &password); err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	id := service.RegisterUser(&username, &password)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   id,
		Token:    username + "|" + password,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	log.Printf("Login name: %s pwd: %s", username, password)
	if id, err := service.QueryUserExisted(&username, &password); err == nil {
		log.Printf("Query user exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    username + "|" + password,
		})
	} else {
		log.Printf("Query user not exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
}

func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")

	ckId, err := service.CheckTokenReturnID(&token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error()},
		})
		return
	}
	// ??????user_id?????????????????????ckId???????????????
	if user_id == "" {
		user, err := service.QueryUserByUserId(ckId)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error()},
			})
			return
		}
		// ?????????????????????????????????????????????
		user.IsFollow = false
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	} else {
		// ????????????user_id???????????????
		id, err := strconv.ParseInt(user_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Unknow ID"},
			})
			return
		}
		user, err := service.QueryUserByUserId(id)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error()},
			})
			return
		}
		// ???????????????????????????????????????user_id
		user.IsFollow = service.IsUserFollowToUser(ckId, id)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	}
}
