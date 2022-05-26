package database

import (
	"time"
)

type User struct {
	UserId        int64          `gorm:"column:user_id;primaryKey"`
	Name          string         `gorm:"column:name"`
	PassWord      string         `gorm:"column:password"`
	Follow        []UserRelation `gorm:"ForeignKey:FollowId"`
	Follower      []UserRelation `gorm:"ForeignKey:FanId"`
	PublishVideos []Video        `gorm:"ForeignKey:Publisher"`
	LikeVideos    []LikeList     `gorm:"ForeignKey:UserLikedId"`
	Comments      []Comment      `gorm:"ForeignKey:CommentUserId"`
}

/***************    table users    ***************

+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| user_id  | bigint       | NO   | PRI | NULL    | auto_increment |
| name     | varchar(256) | YES  |     | NULL    |                |
| password | varchar(256) | YES  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+

**************************************************/

type UserRelation struct {
	FollowId int64 `gorm:"column:follow_id"`
	FanId    int64 `gorm:"column:fan_id"`
}

/***********    table use_relations    ***********

+-----------+--------+------+-----+---------+-------+
| Field     | Type   | Null | Key | Default | Extra |
+-----------+--------+------+-----+---------+-------+
| follow_id | bigint | YES  | MUL | NULL    |       |
| fan_id    | bigint | YES  | MUL | NULL    |       |
+-----------+--------+------+-----+---------+-------+

**************************************************/

type Video struct {
	VideoId    int64      `gorm:"column:video_id;primaryKey"`
	Publisher  int64      `gorm:"column:publisher"`
	Title      string     `gorm:"column:title"`
	PlayUrl    string     `gorm:"column:play_url"`
	CoverUrl   string     `gorm:"column:cover_url"`
	CreateTime int64      `gorm:"column:created_time"`
	LikeUsers  []LikeList `gorm:"ForeignKey:VideoLikedId"`
	Comments   []Comment  `gorm:"ForeignKey:CommentVideoId"`
}

/***************    table videos    **************

+-----------+--------------+------+-----+---------+----------------+
| Field     | Type         | Null | Key | Default | Extra          |
+-----------+--------------+------+-----+---------+----------------+
| video_id  | bigint       | NO   | PRI | NULL    | auto_increment |
| publisher | bigint       | YES  | MUL | NULL    |                |
| play_url  | varchar(256) | YES  |     | NULL    |                |
| cover_url | varchar(256) | YES  |     | NULL    |                |
+-----------+--------------+------+-----+---------+----------------+

**************************************************/

type LikeList struct {
	VideoLikedId int64 `gorm:"column:video_id"`
	UserLikedId  int64 `gorm:"column:user_id"`
}

/***********    table like_lists    ***********

+----------+--------+------+-----+---------+-------+
| Field    | Type   | Null | Key | Default | Extra |
+----------+--------+------+-----+---------+-------+
| video_id | bigint | YES  | MUL | NULL    |       |
| user_id  | bigint | YES  | MUL | NULL    |       |
+----------+--------+------+-----+---------+-------+

**************************************************/

type Comment struct {
	CommentId      int64     `gorm:"column:comment_id;primaryKey"`
	CommentVideoId int64     `gorm:"column:video_id"`
	CommentUserId  int64     `gorm:"column:user_id"`
	Content        string    `gorm:"column:content"`
	CreateAt       time.Time `gorm:"column:create_time"`
}

/***********    table comments    ***********

+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| comment_id  | bigint       | NO   | PRI | NULL    | auto_increment |
| video_id    | bigint       | YES  | MUL | NULL    |                |
| user_id     | bigint       | YES  | MUL | NULL    |                |
| content     | varchar(256) | YES  |     | NULL    |                |
| create_data | datetime     | YES  |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+

**************************************************/
