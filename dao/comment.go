package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID      int64  `gorm:"primary_key"`
	Content string `gorm:"not null"`
	UserID  int64  `gorm:"not null"`
	PostID  int64  `gorm:"not null"`
	Ctime   int64
	Utime   int64
}

type GROMCommentDAO struct {
	db *gorm.DB
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	res := &GROMCommentDAO{
		db: db,
	}
	return res
}

type CommentDAO interface {
	Create(ctx context.Context, comment Comment) (int64, error)
	LIST(ctx context.Context, postId int64, offset int, limit int) ([]Comment, error)
}

func (dao *GROMCommentDAO) Create(ctx context.Context, comment Comment) (int64, error) {
	now := time.Now().UnixMilli()
	comment.Ctime = now
	comment.Utime = now
	err := dao.db.WithContext(ctx).Create(&comment).Error
	return comment.ID, err
}

func (dao *GROMCommentDAO) LIST(ctx context.Context, postId int64, offset int, limit int) ([]Comment, error) {
	var comments []Comment
	result := dao.db.WithContext(ctx).Preload("User").Preload("Post").Where("post_id = ?", postId).Offset(offset).Limit(limit).Find(&comments)
	return comments, result.Error
}
