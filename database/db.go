package database

import (
	"errors"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Dao struct {
}

var (
	db      *gorm.DB
	dao     *Dao
	daoOnce sync.Once
)

func NewDaoInstance() *Dao {
	daoOnce.Do(
		func() {
			dao = &Dao{}
		})
	return dao
}

func Init() {
	log.Printf("Start initializing database")
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:pswd_root@tcp(127.0.0.1:3306)/doushengDB?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                 // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		log.Printf("Failed to connect to database")
		panic(err)
	}

	db.AutoMigrate(&User{}, &UserRelation{}, &Video{}, &Comment{}, &LikeList{})
	// db.AutoMigrate(&LikeList{}, &Comment{}, &Video{}, &UserRelation{}, &User{})
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&UserRelation{})
	// db.AutoMigrate(&Video{})
	// db.AutoMigrate(&Comment{})
	// db.AutoMigrate(&LikeList{})

	sqlDB, _ := db.DB()

	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

	log.Printf("Database initialization succeeded")
}

// 获取视频
func (*Dao) QueryVideos(time int64) (*[]Video, error) {
	var videos []Video
	result := db.Model(&Video{}).Where("created_time >= ?", time).
		Limit(10).Order("created_time desc").Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("no videos")
		result = db.Model(&Video{}).Order("created_time desc").
			Limit(10).Find(&videos)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &[]Video{}, nil
		}
		return &videos, nil
	}
	return &videos, nil
}

// 添加用户
func (*Dao) InsertUser(user *User) error {
	result := db.Model(&User{}).Create(user)
	log.Printf("DB insert user id: %d", user.UserId)
	return result.Error
}

// 通过名称和密码查找用户
func (*Dao) QueryUserByNamePwd(name *string, password *string) (*User, error) {
	var user User
	// log.Printf("QueryUserByName name: %s, password: %s", *name, *password)
	result := db.Where("name = ? AND password = ?", *name, *password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

// 通过名称判断用户是否存在
func (*Dao) QueryUserByName(name *string) (*User, error) {
	var user User
	result := db.Model(&User{}).Where("name = ?", *name).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

// 通过id查找用户
func (*Dao) QueryUserByUserId(userId int64) (*User, error) {
	var user User
	result := db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

// 根据id查找关注数
func (*Dao) CountUserFollowById(userId int64) int64 {
	var count int64
	db.Model(&UserRelation{}).Where("fan_id = ?", userId).Count(&count)
	return count
}

// 根据id查找被关注数
func (*Dao) CountUserFollowerById(userId int64) int64 {
	var count int64
	db.Model(&UserRelation{}).Where("follow_id = ?", userId).Count(&count)
	return count
}

// 添加video
func (*Dao) InsertVideo(video *Video) error {
	result := db.Create(video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 删除video
func (*Dao) QueryVideoByVideoId(videoId int64) *Video {
	var video Video
	db.First(&video, videoId)
	return &video
}

// 根据userId查找发布的video
func (*Dao) QueryVideosByUserId(userId int64) (*[]Video, error) {
	var videos []Video
	result := db.Model(&Video{}).Where("publisher = ?", userId).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &[]Video{}, nil
	}
	return &videos, nil
}

// 根据userId查找喜爱的video
func (*Dao) QueryVideosLikedByUserId(userId int64) (*[]Video, error) {
	var videos []Video
	result := db.Model(&Video{}).
		Joins("inner join like_lists on like_lists.video_id = videos.video_id").
		Where("like_lists.user_id = ?", userId).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &[]Video{}, nil
	}
	return &videos, nil
}

// 根据videoId查找喜爱人数
func (*Dao) CountVideoLikesByVideoId(videoId int64) int64 {
	var count int64
	db.Model(&LikeList{}).Where("video_id = ?", videoId).Count(&count)
	return count
}

// 根据videoId查找评论数
func (*Dao) CountVideoCommentsByVideoId(videoId int64) int64 {
	var count int64
	db.Model(&Comment{}).Where("video_id = ?", videoId).Count(&count)
	return count
}

// 判断userId是否喜爱videoId
func (*Dao) QueryIsUserLikeVideo(userId int64, videoId int64) bool {
	var count int64
	db.Model(&LikeList{}).
		Where("video_id = ? AND user_id = ?", videoId, userId).
		Count(&count)
	return count == 1
}

// 添加关注，fan添加对user的关注
func (*Dao) InsertRelation(fanId int64, userId int64) error {
	result := db.Model(&UserRelation{}).
		Create(&UserRelation{FollowId: userId, FanId: fanId})
	return result.Error
}

// 取消关注，fan取消对user的关注
func (*Dao) DeleteRelation(fanId int64, userId int64) error {
	result := db.Model(&UserRelation{}).
		Where("follow_id = ? AND fan_id = ?", userId, fanId).
		Delete(&UserRelation{})
	return result.Error
}

// 判断user是否关注了to_user
func (*Dao) QueryIsUserRelationToUser(userId int64, to_userId int64) bool {
	var count int64
	db.Model(&UserRelation{}).
		Where("follow_id = ? AND fan_id = ?", to_userId, userId).
		Count(&count)
	return count == 1
}

// 如果action为1，查找user关注的用户列表
// 如果action为2，查找关注user的用户列表
func (*Dao) QueryUserRelationList(userId int64, action int) (*[]User, error) {
	var users []User
	var result *gorm.DB
	if action == 1 {
		log.Printf("get follow list from id: %d ", userId)
		// select (users.user_id, users.name) from users where users.user_id in
		//           (select user_relations.follow_id where user_relations.fan_id = ?), userId
		subQuery := db.Select("follow_id").Where("fan_id = ?", userId).Table("user_relations")
		result = db.Select("user_id, name").Table("users").Where("user_id in (?)", subQuery).Find(&users)
	} else if action == 2 {
		log.Printf("get follower list from id: %d ", userId)
		subQuery := db.Select("fan_id").Where("follow_id = ?", userId).Table("user_relations")
		result = db.Select("user_id, name").Table("users").Where("user_id in (?)", subQuery).Find(&users)
	} else {
		return nil, errors.New("unknow action")
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &[]User{}, nil
	}
	for i := 0; i < len(users); i++ {
		log.Printf(users[i].Name)
	}
	return &users, nil
}

// 添加喜爱，user对video的点赞
func (*Dao) InsertLike(userId int64, videoId int64) error {
	result := db.Model(&LikeList{}).
		Create(&LikeList{VideoLikedId: videoId, UserLikedId: userId})
	return result.Error
}

// 取消喜爱，user取消video的点赞
func (*Dao) DeleteLike(userId int64, videoId int64) error {
	result := db.Model(&LikeList{}).
		Where("video_id = ? AND user_id = ?", videoId, userId).
		Delete(&LikeList{})
	return result.Error
}

// 添加评论，返回评论Id
func (*Dao) InsertComment(userId int64, videoId int64, text *string) (int64, error) {
	comment := Comment{
		CommentVideoId: videoId,
		CommentUserId:  userId,
		Content:        *text,
		CreateAt:       time.Now(),
	}
	result := db.Model(&comment).
		Create(&comment)
	if result.Error != nil {
		return 0, result.Error
	}
	return comment.CommentId, nil
}

// 删除评论
func (*Dao) DeleteComment(commentId int64) error {
	result := db.Delete(&Comment{}, commentId)
	return result.Error
}

// 通过videoId查找评论
func (*Dao) QueryCommentListByVideoId(videoId int64) (*[]Comment, error) {
	var comments []Comment
	result := db.Model(&Comment{}).
		Where("video_id = ?", videoId).
		Order("create_time desc").
		Find(&comments)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &[]Comment{}, nil
	}
	return &comments, nil
}
