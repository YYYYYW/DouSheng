package service

import (
	"DouSheng/database"
	"errors"
	"time"
)

func UserPublishComment(userId int64, videoId int64, text *string) (*Comment, error) {
	exist, user := QueryUserByUserId(userId)
	if !exist {
		return nil, errors.New("User doesn't exist")
	}
	commentId, err := database.NewDaoInstance().InsertComment(userId, videoId, text)
	if err != nil {
		return nil, err
	}
	commentCtr := Comment{
		Id:         commentId,
		User:       *user,
		Content:    *text,
		CreateDate: time.Now().Format("MM-dd"),
	}
	return &commentCtr, nil
}

func UserDeleteComment(commentId int64) error {
	return database.NewDaoInstance().DeleteComment(commentId)
}

func QueryCommentListByVideoId(videoId int64) ([]Comment, error) {
	commentsDB, err := database.NewDaoInstance().QueryCommentListByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	commentsLen := len(commentsDB)
	commentsCtr := make([]Comment, commentsLen)
	for i := 0; i < commentsLen; i++ {
		_, user := QueryUserByUserId(commentsDB[i].CommentUserId)
		commentsCtr[i] = Comment{
			Id:         commentsDB[i].CommentId,
			User:       *user,
			Content:    commentsDB[i].Content,
			CreateDate: commentsDB[i].CreateAt.Format("MM-dd"),
		}
	}
	return commentsCtr, nil
}
