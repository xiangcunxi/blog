package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Comments []Comment
}

type GROMUserDAO struct {
	db *gorm.DB
}

type UserDAO interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
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

func (dao *GROMUserDAO) CreateUser(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Create(&u).Error
	return err
}

func (dao *GROMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.Debug().WithContext(ctx).Where("id = ?", id).First(&u).Error
	fmt.Println("error:", err)
	return u, err
}
