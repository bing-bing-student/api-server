package middleware

import (
	"blog/controller"
	"blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWT 定义中间件, 进行用户权限校验
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		if accessToken != "" {
			// 校验accessToken
			claims, err := utils.ParseAccessToken(accessToken)
			if err == nil {
				// 短token没有过期的情况对请求放行
				c.Set("UserID", claims.UserID)
				c.Next()
				return
			}

			// 如果短token无效或过期, 返回特定的状态码, 让前端去拿长token过来, 去请求/refresh-token接口
			c.JSON(http.StatusUnauthorized, controller.Response{Msg: "access token expired"})
			c.Abort()
			return
		}

		// 如果没有提供任何token, 则返回错误
		c.JSON(http.StatusForbidden, controller.Response{Code: 1002, Msg: "token is required"})
		c.Abort()
	}
}
