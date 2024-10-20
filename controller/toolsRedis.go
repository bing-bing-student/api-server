package controller

import (
	"blog/global"
	"blog/service/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ToolsList struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func ReturnToolsList(c *gin.Context) {
	var tools ToolsList
	var toolsList []ToolsList

	toolsIdList, err := global.RedisSentinel.SMembers(global.CONTEXT, "toolsSet").Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	for _, z := range toolsIdList {
		//查看工具信息是否在redis缓存中
		exist, _ := global.RedisSentinel.Exists(global.CONTEXT, "tools:"+z).Result()
		if exist == 0 {
			//从mysql中获取数据
			id, _ := strconv.ParseUint(z, 10, 64)
			toolsMysql, result := mysql.GetToolsById(id)
			if result != 0 {
				tools.Id = strconv.FormatUint(toolsMysql.ToolID, 10)
				tools.Name = toolsMysql.Describe
				tools.Url = toolsMysql.URL
			}
			mapData := map[string]interface{}{
				"id":   tools.Id,
				"url":  tools.Url,
				"name": tools.Name,
			}
			global.RedisSentinel.HMSet(global.CONTEXT, "tools:"+z, mapData)
		} else {
			//从redis中获取数据
			fields := []string{"id", "name", "url"}
			result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "tools:"+z, fields...).Result()

			tools.Id = result[0].(string)
			tools.Name = result[1].(string)
			tools.Url = result[2].(string)
		}
		toolsList = append(toolsList, tools)
	}
	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"toolsList": toolsList,
	})
}
