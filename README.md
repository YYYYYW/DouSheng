# DouSheng | 青训营极简版抖音项目

## 使用技术

- gin
- gorm
- mysql
- ffmpeg



## 项目结构

整个项目结构分为三层，**数据层**，**逻辑层**，**控制层**，分别对应项目中的**database文件夹**、**service文件夹**、**controller文件夹**。

整体的调用逻辑为：

**前端调用控制层，控制层调用逻辑层，逻辑层调用数据层，数据层从数据库中读写数据。**

静态视频放在**public文件夹**中。

### 数据层

#### gorm数据模型

- 一个用户会有一个id，一个名称，一个密码，有多个关注，多个被关注，多个喜爱的视频，多个发布的视频，多个评论
- 一个视频只有一个视频id，一个发布者，一个标题，一个播放url，一个封面url，一个创造时间，会有多个喜爱用户，多条评论
- 一个评论只有一个评论id，一个视频id，一个用户id，一个内容，一个时间

```go
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

type UserRelation struct {
	FollowId int64 `gorm:"primaryKey;autoIncrement:false;column:follow_id;not null"`
	FanId    int64 `gorm:"primaryKey;autoIncrement:false;column:fan_id;not null"`
}

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

type LikeList struct {
	VideoLikedId int64 `gorm:"primaryKey;autoIncrement:false;column:video_id;not null"`
	UserLikedId  int64 `gorm:"primaryKey;autoIncrement:false;column:user_id;not null"`
}

type Comment struct {
	CommentId      int64     `gorm:"column:comment_id;primaryKey"`
	CommentVideoId int64     `gorm:"index;column:video_id;not null"`
	CommentUserId  int64     `gorm:"column:user_id;not null"`
	Content        string    `gorm:"column:content"`
	CreateAt       time.Time `gorm:"column:create_time;not null"`
}
```

#### MySQL表结构

- **users**表中，**user_id**为主键，**name**为唯一的，并且给**name**添加了索引
- **user_relations**表中，**follow_id**和**fan_id**为复合主键，并且两个都是**users**表中**user_id**的外键
- **videos**表中，**video_id**为主键，**publisher**为**users**表中**user_id**的外键，并添加了索引
- **like_lists**表中，**video_id**和**user_id**为复合主键，**video_id**是**videos**表中**video_id**的外键，**user_id**为**users**表中**user_id**的外键
- **comments**表中，**comment_id**为主键，**video_id**是**videos**表中**video_id**的外键，**user_id**为**users**表中user_id的外键

```
/***************    table users    ***************

+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| user_id  | bigint       | NO   | PRI | NULL    | auto_increment |
| name     | varchar(256) | NO   | UNI | NULL    |                |
| password | varchar(256) | NO   |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+

**************************************************/

/***********    table user_relations    ***********

+-----------+--------+------+-----+---------+-------+
| Field     | Type   | Null | Key | Default | Extra |
+-----------+--------+------+-----+---------+-------+
| follow_id | bigint | NO   | PRI | NULL    |       |
| fan_id    | bigint | NO   | PRI | NULL    |       |
+-----------+--------+------+-----+---------+-------+

**************************************************/

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

/***********    table like_lists    ***********

+----------+--------+------+-----+---------+-------+
| Field    | Type   | Null | Key | Default | Extra |
+----------+--------+------+-----+---------+-------+
| video_id | bigint | NO   | PRI | NULL    |       |
| user_id  | bigint | NO   | PRI | NULL    |       |
+----------+--------+------+-----+---------+-------+

**************************************************/

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
```

#### db.go文件

- Init函数用于连接数据库，并设置数据库连接池参数
- 使用gorm.DB来访问数据库，读写数据
- 使用一个单例对象Dao提供给逻辑层，逻辑层调用Dao的函数



### 逻辑层

逻辑层做这样的一个事情，将从数据层获取到的数据转换成控制层需要的数据结构

#### query_user_info.go

包含一些函数用于查找有关用户的信息



#### query_comment_info.go

包含一些函数用于查找有关评论的信息



#### query_video_info.go

包含一些函数用于查找有关视频的信息



#### common.go

包含逻辑层用于返回给控制层的数据结构



### 控制层

实现了前端需要调用的一些接口

- Feed
- UserInfo
- Register
- Login
- Publish
- PublishList
- FavoriteAction
- FavoriteList
- CommentAction
- CommentList
- RelationAction
- FollowList
- FollowerList



## 一些问题说明

### 关于token（鉴权操作）

- token做的比较简单，并且没有将token存储在数据库。当用户登录或者注册时，会将"name|password"这样一串字符串作为token返回。
- 鉴权操作的流程
  1. 当获得一个token，如果为空，返回errors.New("please login")，否则->2。（Feed流中例外，而是会返回视频列表）
  2. 根据token，先在map中查找是否有该token对应的用户信息，如果有则鉴权成功，否则->3
  3. 根据token得到name和password，在数据库中查找这一项，有则返回用户id，鉴权成功（并将用户信息和token添加到map中），否则鉴权失败



### 关于数据库索引

- 除了所有的主键会自动添加索引外
- 对user表中name列添加了索引，加速通过name查找id
- 对video表中publisher列添加索引，加速查找某个用户发布的视频
- 对comment表中video_id列表添加索引，加速查找某个视频下所有的评论



### 关于视频和封面存放

- 用户上传视频后，视频将会放在public文件夹中
- 视频封面使用ffmpeg截取视频中的第一帧，并将封面图片放在public文件夹下
- 其封面文件名为”视频文件名.jpg“



