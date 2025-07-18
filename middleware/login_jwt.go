package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type LoginJWTMiddleware struct {
	paths []string
}

func NewLoginJWTMiddleware() *LoginJWTMiddleware {
	return &LoginJWTMiddleware{}
}

func (l *LoginJWTMiddleware) IgnorePath(path string) *LoginJWTMiddleware {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("jwt-token")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("Qk1Qb2p6b3h1b1l6b2p6b3h1b1l6b2p6b3h1b1l6b2p6b3h1b1l6b2o="), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
