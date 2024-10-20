package controller

import (
	"blog/global"
	"blog/model"
	"blog/service/mysql"
	"blog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// UserLoginRequest 用户登录的请求结构体
type UserLoginRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type ReturnCodeRequest struct {
	PhoneNumber string `form:"phoneNumber" binding:"required"`
}

type UserLoginByPhoneRequest struct {
	Phone            string `form:"phone" binding:"required,len=11"`
	ShortMessageCode string `form:"shortMessageCode" binding:"required,len=6"`
}

// LoginByToken 用户名密码登录
func LoginByToken(c *gin.Context) {
	// 参数绑定
	var userLogin UserLoginRequest
	if err := c.ShouldBind(&userLogin); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	// 参数校验
	if err := utils.CustomRules("login"); err != nil {
		global.LOGGER.Error("自定义规则错误:", zap.Error(err))
		return
	}

	// 从数据库查询用户信息
	userModel, err := mysql.Login(userLogin.UserName, userLogin.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Code": 1000})
		return
	}

	userId := strconv.FormatUint(userModel.UserID, 10)
	// 返回成功并生成响应 json
	c.JSON(http.StatusOK, gin.H{"Code": 2000, "UserId": userId})
	return
}

// ReturnCode 返回验证码
func ReturnCode(c *gin.Context) {
	// 接收手机号
	var returnCodeRequest *ReturnCodeRequest
	if err := c.ShouldBind(&returnCodeRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	// 校验手机号
	if exist := mysql.IfExistPhoneNumber(returnCodeRequest.PhoneNumber); exist == false {
		c.JSON(http.StatusBadRequest, Response{
			Code: 1003,
			Msg:  "非法手机号",
		})
		return
	}
	// 发送验证码
	err := utils.SendMessage(returnCodeRequest.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 1004,
			Msg:  "验证码生成失败",
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, Response{
		Code: 1005,
		Msg:  "验证码发送成功",
	})
}

// LoginBySMS 短信登陆
func LoginBySMS(c *gin.Context) {
	// 参数绑定
	var userLoginByPhoneRequest *UserLoginByPhoneRequest
	if err := c.ShouldBind(&userLoginByPhoneRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		c.JSON(http.StatusBadRequest, Response{
			Code: 1006,
			Msg:  "无效的手机号",
		})
	}
	// 从缓存中拿到验证码
	code, err := global.RedisSentinel.Get(global.CONTEXT, "code").Result()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: 1004,
			Msg:  "验证码过期",
		})
		return
	}
	// 参数校验
	if userLoginByPhoneRequest.ShortMessageCode != code {
		c.JSON(http.StatusOK, Response{
			Code: 1003,
			Msg:  "手机号或验证码错误",
		})
		return
	}
	// 生成双 token
	userId := mysql.GetUserId(&userLoginByPhoneRequest.Phone)
	claims := new(utils.UserClaims)
	claims.UserID = userId
	doubleToken, err := utils.GenerateToken(claims)
	if err != nil {
		global.LOGGER.Error("生成token错误:", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UserId":       userId,
		"AccessToken":  doubleToken.AccessToken,
		"RefreshToken": doubleToken.RefreshToken,
	})
}

func ManageAboutMe(c *gin.Context) {
	//参数绑定
	type user struct {
		UserID       string `json:"user_id" binding:"required"`
		Name         string `form:"name"`
		Sex          string `form:"sex"`
		Profession   string `form:"profession"`
		Position     string `form:"position"`
		Language     string `form:"language"`
		Domain       string `form:"domain"`
		Introduction string `form:"introduction"`
		Location     string `form:"location"`
		Email        string `form:"email"`
	}
	var userRequestInfo user
	if err := c.ShouldBind(&userRequestInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	//插入数据
	var userInfo model.UserInfo
	userInfo.UserID, _ = strconv.ParseUint(userRequestInfo.UserID, 10, 64)
	userInfo.Name = userRequestInfo.Name
	userInfo.Sex = userRequestInfo.Sex
	userInfo.Email = userRequestInfo.Email
	userInfo.Location = userRequestInfo.Location
	userInfo.Language = userRequestInfo.Language
	userInfo.Profession = userRequestInfo.Profession
	userInfo.Domain = userRequestInfo.Domain
	userInfo.Position = userRequestInfo.Position
	userInfo.Introduction = userRequestInfo.Introduction

	//延迟双删
	err := global.RedisSentinel.Del(global.CONTEXT, "managerInfo").Err()
	result := mysql.ModifyAboutMe(&userInfo)
	time.Sleep(200 * time.Millisecond)
	err = global.RedisSentinel.Del(global.CONTEXT, "managerInfo").Err()

	//返回结果
	if result != 0 && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "提交成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "提交失败",
	})
}
