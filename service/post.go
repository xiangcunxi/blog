package service

import (
	"blog/dao"
	"blog/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("文章参数绑定错误", zap.Error(err))
		return
	}

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
			Msg:  "用户ID错误",
		})
		return
	}
	userId := int64(userIdFloat)

	if req.Id > 0 {
		post, err := p.dao.FindById(ctx, req.Id)
		if err != nil {
			ctx.JSON(http.StatusOK, domain.Result{
				Code: 400,
				Msg:  "文章不存在",
			})
			zap.L().Error("文章不存在", zap.Error(err), zap.Int64("post_id", req.Id))
			return
		}
		if post.Author != userId {
			ctx.JSON(http.StatusOK, domain.Result{
				Code: 400,
				Msg:  "没有修改权限",
			})
			zap.L().Error("没有修改权限", zap.Int64("post_id", req.Id), zap.Int64("user_id", userId))
			return
		}
		err = p.dao.UpdateById(ctx, dao.Post{
			ID:      req.Id,
			Title:   req.Title,
			Content: req.Content,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, domain.Result{
				Code: 500,
				Msg:  "文章更新失败",
			})
			zap.L().Error("文章更新失败", zap.Error(err), zap.Int64("post_id", req.Id))
			return
		}
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 200,
			Msg:  "文章更新成功",
		})
		return
	}

	id, err := p.dao.Create(ctx, dao.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  userId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "文章创建失败",
		})
		zap.L().Error("文章创建失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "文章创建成功",
		Data: id,
	})
}

func (p *PostHandler) Delete(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("参数错误", zap.Error(err), zap.String("param", idstr))
		return
	}
	//获取用户ID
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
			Msg:  "用户ID错误",
		})
		return
	}
	userId := int64(userIdFloat)

	post, err := p.dao.FindById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "文章不存在",
		})
		zap.L().Error("删除文章不存在", zap.Error(err), zap.Int64("post_id", id))
		return
	}
	if post.Author != userId {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "没有删除权限",
		})
		zap.L().Error("没有删除权限", zap.Int64("post_id", id), zap.Int64("user_id", userId))
		return
	}

	err = p.dao.DeleteById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "删除文章失败",
		})
		zap.L().Error("删除文章失败", zap.Error(err), zap.Int64("post_id", id))
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "删除文章成功",
	})
}

func (p *PostHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("参数错误", zap.Error(err), zap.String("param", idstr))
		return
	}
	postList, err := p.dao.FindById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "查询文章详情不存在",
		})
		zap.L().Error("查询文章详情不存在", zap.Error(err), zap.Int64("post_id", id))
		return
	}

	usr, err := p.userDao.FindById(ctx, postList.Author)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "查询文章作者失败",
		})
		zap.L().Error("查询文章作者失败", zap.Error(err), zap.Int64("user_id", postList.Author))
		return
	}

	res := PostVO{
		Id:      postList.ID,
		Title:   postList.Title,
		Content: postList.Content,
		Author:  usr.Username,
		Ctime:   postList.Ctime,
		Utime:   postList.Utime,
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "查询文章详情成功",
		Data: res,
	})
}

func (p *PostHandler) List(ctx *gin.Context) {
	type ListReq struct {
		Offest int `json:"offset"`
		Limit  int `json:"limit"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "参数错误",
		})
		zap.L().Error("获取文章列表参数绑定错误", zap.Error(err))
		return
	}

	userIdInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 401,
			Msg:  "用户未登录",
		})
		return
	}

	userIdFloat, ok := userIdInterface.(float64)
	userId := int64(userIdFloat)
	if !ok {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 400,
			Msg:  "用户ID错误",
		})
		return
	}

	res, err := p.dao.List(ctx, userId, req.Offest, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 500,
			Msg:  "获取文章列表失败",
		})
		zap.L().Error("获取文章列表失败", zap.Error(err))
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
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Author:  authorName,
			Ctime:   post.Ctime,
			Utime:   post.Utime,
		})
	}

	ctx.JSON(http.StatusOK, domain.Result{
		Code: 200,
		Msg:  "获取文章列表成功",
		Data: voList,
	})
}
