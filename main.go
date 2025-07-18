package main

import (
	"blog/dao"
	"blog/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:xiang123@tcp(192.168.29.128:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	dao.Run(db)
	userDao := dao.NewUserDAO(db)

	server := gin.Default()
	u := service.NewUserHandler(userDao)
	u.RegisterRoutes(server)
	server.Run(":8080")
}
