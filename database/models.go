package database

import (
	"time"
)

type User struct {
	UserId        int64          `gorm:"column:user_id;primaryKey"`
	Name          string         `gorm:"index;column:name;unique;not null"`
	PassWord      string         `gorm:"column:password;not null"`
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
| name     | varchar(256) | NO   | UNI | NULL    |                |
| password | varchar(256) | NO   |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+

**************************************************/

type UserRelation struct {
	FollowId int64 `gorm:"primaryKey;autoIncrement:false;column:follow_id;not null"`
	FanId    int64 `gorm:"primaryKey;autoIncrement:false;column:fan_id;not null"`
}

/***********    table use_relations    ***********

+-----------+--------+------+-----+---------+-------+
| Field     | Type   | Null | Key | Default | Extra |
+-----------+--------+------+-----+---------+-------+
| follow_id | bigint | NO   | PRI | NULL    |       |
| fan_id    | bigint | NO   | PRI | NULL    |       |
+-----------+--------+------+-----+---------+-------+

**************************************************/

type Video struct {
	VideoId    int64      `gorm:"column:video_id;primaryKey"`
	Publisher  int64      `gorm:"index;column:publisher;not null"`
	Title      string     `gorm:"column:title"`
	PlayUrl    string     `gorm:"column:play_url;not null"`
	CoverUrl   string     `gorm:"column:cover_url;not null"`
	CreateTime int64      `gorm:"column:created_time;not null"`
	LikeUsers  []LikeList `gorm:"ForeignKey:VideoLikedId"`
	Comments   []Comment  `gorm:"ForeignKey:CommentVideoId"`
}

/***************    table videos    **************

+--------------+--------------+------+-----+---------+----------------+
| Field        | Type         | Null | Key | Default | Extra          |
+--------------+--------------+------+-----+---------+----------------+
| video_id     | bigint       | NO   | PRI | NULL    | auto_increment |
| publisher    | bigint       | NO   | MUL | NULL    |                |
| title        | varchar(256) | YES  |     | NULL    |                |
| play_url     | varchar(256) | NO   |     | NULL    |                |
| cover_url    | varchar(256) | NO   |     | NULL    |                |
| created_time | bigint       | NO   |     | NULL    |                |
+--------------+--------------+------+-----+---------+----------------+

**************************************************/

type LikeList struct {
	VideoLikedId int64 `gorm:"primaryKey;autoIncrement:false;column:video_id;not null"`
	UserLikedId  int64 `gorm:"primaryKey;autoIncrement:false;column:user_id;not null"`
}

/***********    table like_lists    ***********

+----------+--------+------+-----+---------+-------+
| Field    | Type   | Null | Key | Default | Extra |
+----------+--------+------+-----+---------+-------+
| video_id | bigint | NO   | PRI | NULL    |       |
| user_id  | bigint | NO   | PRI | NULL    |       |
+----------+--------+------+-----+---------+-------+

**************************************************/

type Comment struct {
	CommentId      int64     `gorm:"column:comment_id;primaryKey"`
	CommentVideoId int64     `gorm:"index;column:video_id;not null"`
	CommentUserId  int64     `gorm:"column:user_id;not null"`
	Content        string    `gorm:"column:content"`
	CreateAt       time.Time `gorm:"column:create_time;not null"`
}

/***********    table comments    ***********

+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| comment_id  | bigint       | NO   | PRI | NULL    | auto_increment |
| video_id    | bigint       | NO   | MUL | NULL    |                |
| user_id     | bigint       | NO   | MUL | NULL    |                |
| content     | varchar(256) | YES  |     | NULL    |                |
| create_time | datetime     | NO   |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+

**************************************************/
