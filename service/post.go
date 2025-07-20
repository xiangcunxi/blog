package service

import (
	"blog/dao"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	dao dao.PostDAO
}

func NewPostHandler(dao dao.PostDAO) *PostHandler {
	return &PostHandler{dao: dao}
}

func (p *PostHandler) RegisterRoutes(server *gin.Engine) {
	pg := server.Group("/posts")
	pg.POST("/create", p.Create)
}

func (p *PostHandler) Create(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "用户未登录"})
	}

	userID, ok := userIDInterface.(float64)
	if !ok {
		ctx.JSON(500, gin.H{"error": "用户ID错误"})
	}

	id, err := p.dao.Create(ctx, dao.Post{
		Id:       req.Id,
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: int64(userID),
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, gin.H{"msg": "OK", "Data": id})
}
