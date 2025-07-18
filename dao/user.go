package dao

import (
	"blog/domain"
	"context"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

type GROMUserDAO struct {
	db *gorm.DB
}

type UserDAO interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, u domain.User) error
}

func NewUserDAO(db *gorm.DB) UserDAO {
	res := &GROMUserDAO{
		db: db,
	}
	return res
}

func (dao *GROMUserDAO) FindByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	return u, err
}

func (dao *GROMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GROMUserDAO) CreateUser(ctx context.Context, u domain.User) error {
	err := dao.db.WithContext(ctx).Create(&u).Error
	return err
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
}
