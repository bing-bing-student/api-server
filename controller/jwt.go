package controller

import (
	"blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RefreshToken(c *gin.Context) {
	// 获取刷新令牌
	refreshToken := c.GetHeader("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	if refreshToken != "" {
		// 解析刷新令牌
		claims, err := utils.ParseRefreshToken(refreshToken)
		if err == nil {
			// 如果刷新令牌有效, 重新生成访问令牌和刷新令牌
			doubleToken, err := utils.GenerateToken(claims)
			if err != nil {
				// 长token也过期了, 返回403
				c.JSON(http.StatusForbidden, Response{Msg: "unable to generate new tokens"})
				return
			}

			// 将新访问令牌和刷新令牌发送回客户端
			c.JSON(http.StatusOK, gin.H{
				"UserId":        claims.UserID,
				"access_token":  doubleToken.AccessToken,
				"refresh_token": doubleToken.RefreshToken,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Msg: "refreshToken err"})
		return
	}
	c.JSON(http.StatusForbidden, Response{Msg: "token is required"})
	return
}
