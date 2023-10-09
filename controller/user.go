package controller

import (
	"blog/service"
	"blog/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	Response Response
	UserID   interface{} `json:"userID"`
	Token    interface{} `json:"token"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	fmt.Println(username, password)
	// 注册用户到数据库
	userModel, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// 生成对应 token
	tokenString, err := utils.GenerateToken(userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// 返回成功并生成响应 json
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
		UserID:   userModel.UserID,
		Token:    tokenString,
	})
}

func Login(c *gin.Context) {
	var m map[string]string

	body, err := c.GetRawData()
	if err != nil {
		fmt.Println("从Body中获取参数错误")
	}

	// 反序列化
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("反序列化错误")
	}
	username := m["username"]
	password := m["password"]

	// 从数据库查询用户信息
	userModel, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户名或密码错误"})
		return
	}
	// 生成对应 token
	tokenString, err := utils.GenerateToken(userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	// 返回成功并生成响应 json
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserID: userModel.UserID,
		Token:  tokenString,
	},
	)
}
