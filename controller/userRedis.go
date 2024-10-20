package controller

import (
	"blog/global"
	"blog/service/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ManagerInfo struct {
	Name         string `json:"name"`
	Sex          string `json:"sex"`
	Profession   string `json:"profession"`
	Position     string `json:"position"`
	Language     string `json:"language"`
	Domain       string `json:"domain"`
	Introduction string `json:"introduction"`
	Location     string `json:"location"`
	Email        string `json:"email"`
}

func ReturnAboutMeInfo(c *gin.Context) {
	exist, _ := global.RedisSentinel.Exists(global.CONTEXT, "managerInfo").Result()
	var managerInfo ManagerInfo
	if exist == 0 {
		//从mysql中获取数据
		rootInfo, result := mysql.GetManagerInfo()
		if result == false {
			c.JSON(http.StatusOK, gin.H{
				"msg": "没有相关信息",
			})
			return
		}
		//将数据缓存到redis
		mapData := map[string]interface{}{
			"name":         rootInfo.Name,
			"sex":          rootInfo.Sex,
			"profession":   rootInfo.Profession,
			"position":     rootInfo.Position,
			"language":     rootInfo.Language,
			"domain":       rootInfo.Domain,
			"introduction": rootInfo.Introduction,
			"location":     rootInfo.Location,
			"email":        rootInfo.Email,
		}
		global.RedisSentinel.HMSet(global.CONTEXT, "managerInfo", mapData)
		managerInfo.Name = rootInfo.Name
		managerInfo.Sex = rootInfo.Sex
		managerInfo.Profession = rootInfo.Profession
		managerInfo.Position = rootInfo.Position
		managerInfo.Language = rootInfo.Language
		managerInfo.Domain = rootInfo.Domain
		managerInfo.Introduction = rootInfo.Introduction
		managerInfo.Location = rootInfo.Location
		managerInfo.Email = rootInfo.Email
	} else {
		//从redis中获取数据
		fields := []string{"name", "sex", "profession", "position", "language", "domain", "introduction", "location", "email"}
		res, _ := global.RedisSentinel.HMGet(global.CONTEXT, "managerInfo", fields...).Result()
		managerInfo.Name = res[0].(string)
		managerInfo.Sex = res[1].(string)
		managerInfo.Profession = res[2].(string)
		managerInfo.Position = res[3].(string)
		managerInfo.Language = res[4].(string)
		managerInfo.Domain = res[5].(string)
		managerInfo.Introduction = res[6].(string)
		managerInfo.Location = res[7].(string)
		managerInfo.Email = res[8].(string)
	}
	c.JSON(http.StatusOK, gin.H{
		"about": managerInfo,
	})
}
