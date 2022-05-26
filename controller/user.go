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

	if exist, _ := service.QueryUserExisted(&username, &password); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	id := service.RegisterUser(&username, &password)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   id,
		Token:    username + password,
	})

	// token := username + password

	// if _, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	// 	})
	// } else {
	// 	atomic.AddInt64(&userIdSequence, 1)
	// 	newUser := User{
	// 		Id:   userIdSequence,
	// 		Name: username,
	// 	}
	// 	usersLoginInfo[token] = newUser
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   userIdSequence,
	// 		Token:    username + password,
	// 	})
	// }
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if exist, id := service.QueryUserExisted(&username, &password); exist {
		log.Printf("Query user exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    username + password,
		})
	} else {
		log.Printf("Query user not exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist!"},
		})
	}
	// token := username + password

	// if user, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   user.Id,
	// 		Token:    token,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	// 	})
	// }
}

func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")

	// 表示这是第一次注册时登录的用户
	if user_id == "" {
		mapexist, _ := service.QueryUserIdByToken(&token)
		if !mapexist {
			log.Printf("User info not find in map")
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
	} else {
		// 表示这是后续查询的用户
		id, err := strconv.ParseInt(user_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Unknow ID"},
			})
			return
		}
		ret := service.CheckToken(id, &token)
		if ret == -1 {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			return
		} else if ret == -2 {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "authentication failed"},
			})
			return
		}
		_, user := service.QueryUserByUserId(id)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	}

	// if user, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserResponse{
	// 		Response: Response{StatusCode: 0},
	// 		User:     user,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, UserResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	// 	})
	// }
}
