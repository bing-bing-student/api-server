package controller

import (
	"blog/global"
	"blog/model"
	"blog/service/mysql"
	"blog/service/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// GetTools 返回所有工具信息
func GetTools(c *gin.Context) {
	type tools struct {
		ToolsId  string `json:"tools_id"`
		Describe string `json:"describe"`
		Url      string `json:"url"`
	}
	toolsArray, result := mysql.QueryAllTools()
	if result == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "工具箱中没有任何信息",
		})
	}
	var toolsResponse []tools
	for _, toolsInfo := range toolsArray {
		id := strconv.FormatUint(toolsInfo.ToolID, 10)
		toolsResponse = append(toolsResponse, tools{
			ToolsId:  id,
			Describe: toolsInfo.Describe,
			Url:      toolsInfo.URL,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"toolsArray": toolsResponse,
	})
}

// AddTools 添加工具
func AddTools(c *gin.Context) {
	//参数绑定
	type tools struct {
		Url      string `form:"url" binding:"required"`
		Describe string `form:"describe" binding:"required"`
	}
	var toolsRequestInfo tools
	if err := c.ShouldBind(&toolsRequestInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	var toolsInfo model.Tools
	toolsInfo.ToolID, _ = global.IdGenerator.NextID()
	toolsInfo.URL = toolsRequestInfo.Url
	toolsInfo.Describe = toolsRequestInfo.Describe
	result := mysql.CreateTools(&toolsInfo)

	//在redis中添加一条记录
	err := redis.AddToolsIdToRedis(toolsInfo.ToolID)
	if result == 0 && err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "添加成功",
	})
}

// UpdateTools 修改工具信息
func UpdateTools(c *gin.Context) {
	//参数绑定
	type tools struct {
		ToolsId  string `form:"tools_id" json:"tools_id" binding:"required"`
		Url      string `form:"url" binding:"required"`
		Describe string `form:"describe" binding:"required"`
	}
	var toolsRequestInfo tools
	if err := c.ShouldBind(&toolsRequestInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	//数据转换、更新数据
	var toolsInfo model.Tools
	toolsInfo.ToolID, _ = strconv.ParseUint(toolsRequestInfo.ToolsId, 10, 64)
	toolsInfo.URL = toolsRequestInfo.Url
	toolsInfo.Describe = toolsRequestInfo.Describe

	//延迟双删
	err := global.RedisSentinel.Del(global.CONTEXT, "tools:"+toolsRequestInfo.ToolsId).Err()
	result := mysql.UpdateTools(&toolsInfo)
	time.Sleep(300 * time.Millisecond)
	err = global.RedisSentinel.Del(global.CONTEXT, "tools:"+toolsRequestInfo.ToolsId).Err()

	//返回结果
	if result != 0 && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新失败",
	})
	return
}

// DeleteTools 删除工具
func DeleteTools(c *gin.Context) {
	//参数绑定
	type tools struct {
		ToolsId string `form:"tools_id" binding:"required"`
	}
	var toolsRequestInfo tools
	if err := c.ShouldBind(&toolsRequestInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	//数据转换、更新数据
	toolsId, _ := strconv.ParseUint(toolsRequestInfo.ToolsId, 10, 64)

	//延迟双删
	err := global.RedisSentinel.SRem(global.CONTEXT, "toolsSet", toolsRequestInfo.ToolsId).Err()
	err = global.RedisSentinel.Del(global.CONTEXT, "tools:"+toolsRequestInfo.ToolsId).Err()
	result := mysql.DeleteTools(toolsId)
	time.Sleep(300 * time.Millisecond)
	err = global.RedisSentinel.Del(global.CONTEXT, "tools:"+toolsRequestInfo.ToolsId).Err()

	//返回结果
	if result != 0 && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除失败",
	})
}
