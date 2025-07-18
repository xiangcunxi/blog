package service

import (
	"blog/dao"
	"blog/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserHandler struct {
	dao dao.UserDAO
}

func NewUserHandler(dao dao.UserDAO) *UserHandler {
	return &UserHandler{dao: dao}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := u.dao.FindByUsername(c, req.Username)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户已被注册"})
		return
	}
	_, err = u.dao.FindByEmail(c, req.Email)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮箱已被注册"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = u.dao.CreateUser(c, domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}

	user, err := u.dao.FindByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("Qk1Qb2p6b3h1b1l6b2p6b3h1b1l6b2p6b3h1b1l6b2p6b3h1b1l6b2o="))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	ctx.Header("Authorization", "Bearer "+tokenString)

	ctx.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}
