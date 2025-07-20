package main

import (
	"blog/dao"
	"blog/middleware"
	"blog/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:xiang123@tcp(192.168.29.128:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	dao.InitDB(db)
	userDao := dao.NewUserDAO(db)
	postDao := dao.NewPostDAO(db)

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//不加这个前端拿不到
		ExposeHeaders:    []string{"jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	server.Use(middleware.NewLoginJWTMiddleware().
		IgnorePath("/user/login").
		IgnorePath("/user/register").Build())

	u := service.NewUserHandler(userDao)
	u.RegisterRoutes(server)

	p := service.NewPostHandler(postDao)
	p.RegisterRoutes(server)

	server.Run(":8080")
}
