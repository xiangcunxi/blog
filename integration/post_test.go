package integration

import (
	"blog/dao"
	"blog/service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *PostTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE posts")
}

func (s *PostTestSuite) SetupSuite() {
	s.server = gin.Default()
	db, err := gorm.Open(mysql.Open("root:xiang123@tcp(192.168.29.128:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	s.db = db
	postDao := dao.NewPostDAO(s.db)
	postHdl := service.NewPostHandler(postDao)
	postHdl.RegisterRoutes(s.server)

}

func (s *PostTestSuite) TestCreate() {
	t := s.T()
	testCases := []struct {
		name string

		//准备数据
		before func(t *testing.T)
		//验证数据
		after func(t *testing.T)

		//预期输入
		post Post

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "编辑成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var post dao.Post
				err := s.db.Where("id =?", 1).First(&post).Error
				assert.NoError(t, err)
				assert.True(t, post.Ctime > 0)
				assert.True(t, post.Utime > 0)
				post.Ctime = 0
				post.Utime = 0
				assert.Equal(t, dao.Post{
					Id:      1,
					Title:   "test title",
					Content: "test content",
					Author:  0,
				}, post)
			},
			post: Post{
				Title:   "test title",
				Content: "test content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.post)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/posts/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func TestPost(t *testing.T) {
	suite.Run(t, &PostTestSuite{})
}

type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
