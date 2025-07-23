package service

import (
	"blog/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostHandler struct {
	dao     dao.PostDAO
	userDao dao.UserDAO
}

type PostVO struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Ctime   int64  `json:"ctime"`
	Utime   int64  `json:"utime"`
}

func NewPostHandler(dao dao.PostDAO, userDao dao.UserDAO) *PostHandler {
	return &PostHandler{dao: dao, userDao: userDao}
}

func (p *PostHandler) RegisterRoutes(server *gin.Engine) {
	pg := server.Group("/posts")
	pg.POST("/edit", p.Edit)
	pg.DELETE("/delete/:id", p.Delete)
	pg.GET("/detail/:id", p.Detail)
	pg.POST("/list", p.List)
}

func (p *PostHandler) Edit(ctx *gin.Context) {
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

	userIdInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "用户未登录"})
		return
	}

	userIdFloat, ok := userIdInterface.(float64)
	if !ok {
		ctx.JSON(500, gin.H{"error": "用户ID错误"})
		return
	}
	userId := int64(userIdFloat)

	if req.Id > 0 {
		post, err := p.dao.FindById(ctx, req.Id)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if post.Author != userId {
			ctx.JSON(403, gin.H{"error": "没有修改权限"})
			return
		}
		err = p.dao.UpdateById(ctx, dao.Post{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
		})
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"msg": "修改成功"})
		return
	}

	id, err := p.dao.Create(ctx, dao.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  userId,
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"msg": "创建成功", "id": id})
}

func (p *PostHandler) Delete(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	userIdInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "用户未登录"})
		return
	}
	userIdFloat, ok := userIdInterface.(float64)
	if !ok {
		ctx.JSON(500, gin.H{"error": "用户ID错误"})
		return
	}
	userId := int64(userIdFloat)

	post, err := p.dao.FindById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if post.Author != userId {
		ctx.JSON(403, gin.H{"error": "没有删除权限"})
		return
	}

	err = p.dao.DeleteById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"msg": "删除成功"})
}

func (p *PostHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	postList, err := p.dao.FindById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	usr, err := p.userDao.FindById(ctx, postList.Author)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := PostVO{
		Id:      postList.Id,
		Title:   postList.Title,
		Content: postList.Content,
		Author:  usr.Username,
		Ctime:   postList.Ctime,
		Utime:   postList.Utime,
	}
	ctx.JSON(200, gin.H{"data": res})
}

func (p *PostHandler) List(ctx *gin.Context) {
	type ListReq struct {
		Offest int `json:"offset"`
		Limit  int `json:"limit"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userIdInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "用户未登录"})
		return
	}

	userIdFloat, ok := userIdInterface.(float64)
	userId := int64(userIdFloat)
	if !ok {
		ctx.JSON(500, gin.H{"error": "用户ID错误"})
		return
	}

	res, err := p.dao.List(ctx, userId, req.Offest, req.Limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var voList []PostVO
	for _, post := range res {
		usr, err := p.userDao.FindById(ctx, post.Author)
		authorName := ""
		if err == nil {
			authorName = usr.Username
		}
		voList = append(voList, PostVO{
			Id:      post.Id,
			Title:   post.Title,
			Content: post.Content,
			Author:  authorName,
			Ctime:   post.Ctime,
			Utime:   post.Utime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": voList,
	})
}
