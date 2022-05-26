package service

import (
	"DouSheng/database"
	"errors"
	"time"
)

func UserPublishComment(userId int64, videoId int64, text *string) (*Comment, error) {
	user, err := QueryUserByUserId(userId)
	if err != nil {
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
		CreateDate: time.Now().Format("01-02"),
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
	commentsLen := len(*commentsDB)
	commentsCtr := make([]Comment, commentsLen)
	for i := 0; i < commentsLen; i++ {
		// TODO 可以联查
		user, _ := QueryUserByUserId((*commentsDB)[i].CommentUserId)
		commentsCtr[i] = Comment{
			Id:         (*commentsDB)[i].CommentId,
			User:       *user,
			Content:    (*commentsDB)[i].Content,
			CreateDate: (*commentsDB)[i].CreateAt.Format("01-02"),
		}
	}
	return commentsCtr, nil
}
