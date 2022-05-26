package service

import (
	"DouSheng/database"
	"log"
)

var usersLoginInfo = make(map[string]User)

// 使用Token在map中查找是否存在，找到了返回id
func QueryUserIdByToken(token *string) (bool, int64) {
	if val, exist := usersLoginInfo[*token]; exist {
		return true, val.Id
	}
	return false, -1
}

// 通过用户ID获取用户信息
func QueryUserByUserId(userId int64) (bool, *User) {
	// TODO 可以尝试先根据token直接在map里面找（需要判断是否更新了），不用每次都访问数据库
	log.Printf("Query user by id: %d", userId)
	userDB, err := database.NewDaoInstance().QueryUserByUserId(userId)
	if err != nil {
		return false, nil
	}
	followCount := database.NewDaoInstance().CountUserFollowById(userId)
	followerCount := database.NewDaoInstance().CountUserFollowerById(userId)
	userCtr := User{
		Id:            userDB.UserId,
		Name:          userDB.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
	}
	return true, &userCtr
}

// 根据name，password查询是否存在user，如果存在返回user_id
func QueryUserExisted(name *string, password *string) (bool, int64) {
	token := *name + *password
	if val, exist := usersLoginInfo[token]; exist {
		return true, val.Id
	}
	userDB, err := database.NewDaoInstance().QueryUserByName(name, password)
	if err != nil {
		return false, 0
	}
	userCtr := User{
		Id:   userDB.UserId,
		Name: userDB.Name,
	}
	usersLoginInfo[token] = userCtr
	return true, userDB.UserId
}

// 根据token判断用户是否存在
func QueryUserExistedByToken(token *string) (bool, int64) {
	return QueryUserIdByToken(token)
}

// 根据userId判断用户是否存在
func QueryUserExistedById(userId int64) (*database.User, bool) {
	userDB, err := database.NewDaoInstance().QueryUserByUserId(userId)
	return userDB, err == nil
}

// 用户注册
func RegisterUser(name *string, password *string) int64 {
	token := *name + *password
	newDBUser := database.User{
		Name:     *name,
		PassWord: *password,
	}
	database.NewDaoInstance().InsertUser(&newDBUser)
	log.Printf("RegisterUser name: %s, password: %s, id: %d",
		newDBUser.Name, newDBUser.PassWord, newDBUser.UserId)
	newCtrUser := User{
		Id:   newDBUser.UserId,
		Name: *name,
	}
	usersLoginInfo[token] = newCtrUser
	return newDBUser.UserId
}

// user关注to_user
func UserRelationToUser(userId int64, to_userId int64) error {
	return database.NewDaoInstance().InsertRelation(userId, to_userId)
}

// user取消关注to_user
func UserUnRelationToUser(userId int64, to_userId int64) error {
	return database.NewDaoInstance().DeleteRelation(userId, to_userId)
}

// 查询user的关注列表
func QueryUserFollowList(userId int64) ([]User, error) {
	usersDB, err := database.NewDaoInstance().QueryUserRelationList(userId, 1)
	if err != nil {
		return nil, err
	}
	usersLen := len(usersDB)
	usersCtr := make([]User, usersLen)
	for i := 0; i < usersLen; i++ {
		followerCount := database.NewDaoInstance().CountUserFollowerById(userId)
		usersCtr[i] = User{
			Id:            usersDB[i].UserId,
			Name:          usersDB[i].Name,
			FollowCount:   int64(usersLen),
			FollowerCount: followerCount,
			IsFollow:      true,
		}
	}
	return usersCtr, nil
}

// 查询关注user的列表
func QueryUserFollowerList(userId int64) ([]User, error) {
	usersDB, err := database.NewDaoInstance().QueryUserRelationList(userId, 1)
	if err != nil {
		return nil, err
	}
	usersLen := len(usersDB)
	usersCtr := make([]User, usersLen)
	for i := 0; i < usersLen; i++ {
		followCount := database.NewDaoInstance().CountUserFollowById(userId)
		isFollow := database.NewDaoInstance().QueryIsUserRelationToUser(userId, usersDB[i].UserId)
		usersCtr[i] = User{
			Id:            usersDB[i].UserId,
			Name:          usersDB[i].Name,
			FollowCount:   followCount,
			FollowerCount: int64(usersLen),
			IsFollow:      isFollow,
		}
	}
	return usersCtr, nil
}

// 检查用户Token是否有效，返回为1表示有效，返回-1表示用户不存在，返回-2表示鉴权失败
func CheckToken(userId int64, token *string) int64 {
	exist, id := QueryUserExistedByToken(token)
	if exist {
		if id != userId {
			return -2
		}
		return 1
	}
	if userDB, exist := QueryUserExistedById(userId); exist {
		token := userDB.Name + userDB.PassWord
		usersLoginInfo[token] = User{
			Id:   userDB.UserId,
			Name: userDB.Name,
		}
		return 1
	}
	return -1
}
