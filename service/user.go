package service

import (
	"blog/dao"
	"blog/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		c.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("用户注册绑定参数失败", zap.Error(err))
		return
	}
	_, err := u.dao.FindByUsername(c, req.Username)
	if err == nil {
		c.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "用户名已被注册",
		})
		return
	}
	_, err = u.dao.FindByEmail(c, req.Email)
	if err == nil {
		c.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "邮箱已被注册",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "注册失败",
		})
		zap.L().Error("用户注册密码加密失败", zap.Error(err))
		return
	}

	err = u.dao.CreateUser(c, dao.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "注册失败",
		})
		zap.L().Error("用户注册失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "注册成功",
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("用户登录绑定参数失败", zap.Error(err))
		return
	}
	user, err := u.dao.FindByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "用户不存在",
		})
		zap.L().Info("用户不存在", zap.Error(err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "密码错误",
		})
		zap.L().Info("用户密码错误", zap.Error(err))
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
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "登录失败",
		})
		zap.L().Error("用户登录生成token失败", zap.Error(err))
		return
	}
	ctx.Header("Authorization", "Bearer "+tokenString)

	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "登录成功",
	})
}
