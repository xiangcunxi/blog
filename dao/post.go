package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	Id      int64  `gorm:"primarykey, autoincrement"`
	Title   string `gorm:"type=VARCHAR(1024),not null"`
	Content string `gorm:"type=BLOB, not null"`
	Author  int64  `gorm:"index=pid_ctime"`
	Ctime   int64  `gorm:"index=pid_ctime"`
	Utime   int64
}

type GROMPostDAO struct {
	db *gorm.DB
}

func NewPostDAO(db *gorm.DB) PostDAO {
	res := &GROMPostDAO{
		db: db,
	}
	return res
}

type PostDAO interface {
	Create(ctx context.Context, post Post) (int64, error)
	UpdateById(ctx context.Context, post Post) error
	FindById(ctx context.Context, postId int64) (Post, error)
	DeleteById(ctx context.Context, postId int64) error
	List(ctx context.Context, userId int64, offset int, limit int) ([]Post, error)
}

func (dao *GROMPostDAO) Create(ctx context.Context, post Post) (int64, error) {
	now := time.Now().UnixMilli()
	post.Ctime = now
	post.Utime = now
	err := dao.db.WithContext(ctx).Create(&post).Error
	return post.Id, err
}

func (dao *GROMPostDAO) UpdateById(ctx context.Context, post Post) error {
	now := time.Now().UnixMilli()
	post.Utime = now
	res := dao.db.WithContext(ctx).Model(&post).Where("id = ?", post.Id).
		Updates(map[string]any{
			"title":   post.Title,
			"content": post.Content,
			"utime":   post.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能创作者非法 id %d, author %d", post.Id, post.Author)
	}
	return res.Error
}

func (dao *GROMPostDAO) FindById(ctx context.Context, postId int64) (Post, error) {
	var p Post
	err := dao.db.WithContext(ctx).Where("id=?", postId).First(&p).Error
	return p, err
}

func (dao *GROMPostDAO) DeleteById(ctx context.Context, postId int64) error {
	var p Post
	err := dao.db.WithContext(ctx).Delete(&p, postId).Error
	return err
}

func (dao *GROMPostDAO) List(ctx context.Context, userId int64, offset int, limit int) ([]Post, error) {
	var posts []Post
	err := dao.db.WithContext(ctx).Offset(offset).Limit(limit).Order("utime desc").Find(&posts).Error
	return posts, err
}
