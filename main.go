package main

import (
	"blog/dao"
	"blog/middleware"
	"blog/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	initLogger()
	db, err := gorm.Open(mysql.Open("root:xiang123@tcp(192.168.29.128:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		zap.L().Error("数据库连接失败", zap.Error(err))
		panic(err)
	}
	dao.InitDB(db)
	userDao := dao.NewUserDAO(db)
	postDao := dao.NewPostDAO(db)
	commentDao := dao.NewCommentDAO(db)

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
		IgnorePath("/user/signup").Build())

	u := service.NewUserHandler(userDao)
	u.RegisterRoutes(server)

	p := service.NewPostHandler(postDao, userDao)
	p.RegisterRoutes(server)

	c := service.NewCommentHandler(commentDao, userDao, postDao)
	c.RegisterRoutes(server)

	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
