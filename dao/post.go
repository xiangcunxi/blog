package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	Id       int64  `gorm:"primarykey, autoincrement"`
	Title    string `gorm:"type=VARCHAR(1024),not null"`
	Content  string `gorm:"type=BLOB, not null"`
	UserID   uint
	AuthorId int64 `gorm:"index=pid_ctime"`
	Ctime    int64 `gorm:"index=pid_ctime"`
	Utime    int64
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
}

func (p *GROMPostDAO) Create(ctx context.Context, post Post) (int64, error) {
	now := time.Now().UnixMilli()
	post.Ctime = now
	post.Utime = now
	err := p.db.WithContext(ctx).Create(&post).Error
	return post.Id, err
}
