package service

import (
	"DouSheng/database"
	"errors"
	"log"
	"strings"
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
func QueryUserByUserId(userId int64) (*User, error) {
	// log.Printf("Query user by id: %d", userId)
	userDB, err := database.NewDaoInstance().QueryUserByUserId(userId)
	if err != nil {
		return nil, err
	}
	followCount := database.NewDaoInstance().CountUserFollowById(userId)
	followerCount := database.NewDaoInstance().CountUserFollowerById(userId)
	userCtr := User{
		Id:            userDB.UserId,
		Name:          userDB.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
	}
	return &userCtr, nil
}

// 根据name，password查询是否存在user，如果存在返回user_id
func QueryUserExisted(name *string, password *string) (int64, error) {
	token := *name + "|" + *password
	if val, exist := usersLoginInfo[token]; exist {
		return val.Id, nil
	}
	userDB, err := database.NewDaoInstance().QueryUserByName(name)
	if err != nil {
		return 0, errors.New("User doesn't exist")
	}
	if *password != userDB.PassWord {
		return 0, errors.New("password error")
	}
	userCtr := User{
		Id:   userDB.UserId,
		Name: userDB.Name,
	}
	usersLoginInfo[token] = userCtr
	return userDB.UserId, nil
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

// 判断user是否关注了to_user
func IsUserFollowToUser(userId int64, to_userId int64) bool {
	return database.NewDaoInstance().QueryIsUserRelationToUser(userId, to_userId)
}

// 查询user的关注列表
func QueryUserFollowList(userId int64) ([]User, error) {
	usersDB, err := database.NewDaoInstance().QueryUserRelationList(userId, 1)
	if err != nil {
		return nil, err
	}
	usersLen := len(*usersDB)
	usersCtr := make([]User, usersLen)
	for i := 0; i < usersLen; i++ {
		usersCtr[i] = User{
			Id:   (*usersDB)[i].UserId,
			Name: (*usersDB)[i].Name,
		}
	}
	return usersCtr, nil
}

// 查询关注user的列表
func QueryUserFollowerList(userId int64) ([]User, error) {
	usersDB, err := database.NewDaoInstance().QueryUserRelationList(userId, 2)
	if err != nil {
		return nil, err
	}
	usersLen := len(*usersDB)
	usersCtr := make([]User, usersLen)
	for i := 0; i < usersLen; i++ {
		usersCtr[i] = User{
			Id:   (*usersDB)[i].UserId,
			Name: (*usersDB)[i].Name,
		}
	}
	return usersCtr, nil
}

// 检查用户Token是否有效，有效时返回用户ID，无效时返回err
func CheckTokenReturnID(token *string) (int64, error) {
	if *token == "" {
		return -1, errors.New("please login")
	}
	if exist, id := QueryUserIdByToken(token); exist {
		return id, nil
	}
	sp := strings.Index(*token, "|")
	name := (*token)[:sp]
	password := (*token)[sp+1:]
	userDB, err := database.NewDaoInstance().QueryUserByNamePwd(&name, &password)
	if err != nil {
		return -1, err
	}
	usersLoginInfo[*token] = User{
		Id:   userDB.UserId,
		Name: name,
	}
	return userDB.UserId, nil
}
