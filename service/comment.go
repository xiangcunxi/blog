package service

import (
	"blog/dao"
	"blog/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CommentHandler struct {
	dao     dao.CommentDAO
	userDAO dao.UserDAO
	postDAO dao.PostDAO
}

func NewCommentHandler(dao dao.CommentDAO, userDAO dao.UserDAO, postDAO dao.PostDAO) *CommentHandler {
	return &CommentHandler{dao: dao, userDAO: userDAO, postDAO: postDAO}
}

func (c *CommentHandler) RegisterRoutes(server *gin.Engine) {
	cg := server.Group("/comments")
	cg.POST("/edit", c.Create)
	cg.POST("/list", c.List)
}

func (c *CommentHandler) Create(ctx *gin.Context) {
	type CommentReq struct {
		ID      int64  `json:"id"`
		PostID  int64  `json:"postId"`
		Content string `json:"content"`
	}
	var req CommentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数绑定错误",
		})
		zap.L().Error("创建评论参数绑定错误", zap.Error(err))
		return
	}

	// 获取用户ID
	userIdInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 401,
			Msg:  "用户未登录",
		})
		return
	}
	userIdFloat, ok := userIdInterface.(float64)
	if !ok {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "用户ID类型错误",
		})
		return
	}
	userId := int64(userIdFloat)

	//检查文章是否存在
	_, err := c.postDAO.FindById(ctx, req.PostID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "评论文章不存在",
		})
		zap.L().Error("评论文章不存在", zap.Error(err))
		return
	}

	id, err := c.dao.Create(ctx, dao.Comment{
		ID:      req.ID,
		UserID:  userId,
		PostID:  req.PostID,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "创建评论失败",
		})
		zap.L().Error("创建评论失败", zap.Error(err))
		return
	}
	ctx.JSON(200, domain.Result{
		Code: 200,
		Msg:  "创建评论成功",
		Data: id,
	})
}

func (c *CommentHandler) List(ctx *gin.Context) {
	type ListReq struct {
		PostID int64 `json:"postId"`
		Offest int   `json:"offset"`
		Limit  int   `json:"limit"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数绑定错误",
		})
		zap.L().Error("获取评论列表参数绑定错误", zap.Error(err))
		return
	}
	comments, err := c.dao.LIST(ctx, req.PostID, req.Offest, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "获取评论列表失败",
		})
		zap.L().Error("获取评论列表失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "获取评论列表成功",
		Data: comments,
	})
}
