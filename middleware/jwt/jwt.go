package jwt

import (
	"decept-defense/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = app.SUCCESS
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			code = app.ErrorAuth
		} else {
			extractedToken := strings.Split(token, "Bearer ")
			if len(extractedToken) == 2 {
				token = strings.TrimSpace(extractedToken[1])
				user, err := app.ParseToken(token)
				c.Set("currentUser", user)
				if err != nil {
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						code = app.ErrorAuth
					default:
						code = app.ErrorAuth
					}
				}
			} else {
				code = app.ErrorAuth
			}
		}
		if code != app.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  app.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
