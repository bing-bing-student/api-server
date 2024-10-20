package controller

import (
	"blog/global"
	"blog/model"
	"blog/service/mysql"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"html"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BlogInfo struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Describe   string `json:"describe"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Views      int    `json:"views"`
	Label      string `json:"label"`
}

type BlogCardIndex struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Views      int    `json:"views"`
	Label      string `json:"label"`
}

type BlogCardCategory struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Views      int    `json:"views"`
}

type Article struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Views      int    `json:"views"`
	Label      string `json:"label"`
}

type BlogRequest struct {
	Id string `form:"id" binding:"required,lte=19"`
}

type SearchRequest struct {
	Question string `json:"q" form:"q" binding:"required,lte=20"`
}

type LabelName struct {
	Label string `json:"label" form:"label" binding:"required"`
}

type Label struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"articleCount"`
}

type HotTops struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Views string `json:"views"`
}

type TimeLineArticle struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"createTime"`
}

type YearGroup struct {
	Timestamp string            `json:"timestamp"`
	Articles  []TimeLineArticle `json:"articles"`
}

type Page struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=9"`
}

// GetBlogListOnIndex 返回首页的博客列表
func GetBlogListOnIndex(c *gin.Context) {
	// 参数绑定
	var pageRequest Page
	if err := c.ShouldBindQuery(&pageRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	// 从列表中得到博客id
	page := pageRequest.Page
	pageSize := pageRequest.PageSize

	start := (page - 1) * pageSize
	stop := start + pageSize - 1
	blogIdList, err := global.RedisSentinel.ZRevRangeWithScores(global.CONTEXT, "blogIdList", int64(start), int64(stop)).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	// 遍历blogList
	var blogList []BlogCardIndex
	for _, z := range blogIdList {
		var blog BlogCardIndex
		//查看文章是否在redis缓存中
		exist, _ := global.RedisSentinel.Exists(global.CONTEXT, "article:"+z.Member.(string)).Result()
		if exist == 0 {
			//从mysql中获取文章数据
			var blogInfo *model.Article
			var effective bool

			id, _ := strconv.ParseUint(z.Member.(string), 10, 64)
			if effective, blogInfo = mysql.GetBlogInfoById(id); effective == true {
				label := mysql.GetLabelNameById(blogInfo.LabelID)
				blog.Id = z.Member.(string)
				blog.Title = blogInfo.BlogTitle
				blog.Content = blogInfo.ShowContext
				blog.CreateTime = blogInfo.CreatedTime.Format("2006-01-02")
				blog.Views = int(blogInfo.ViewNum)
				blog.Label = label
			}
			//把文章的所有数据缓存到redis
			mapData := map[string]interface{}{
				"blogId":     blogInfo.ArticleID,
				"title":      blogInfo.BlogTitle,
				"context":    blogInfo.Context,
				"describe":   blogInfo.ShowContext,
				"views":      blogInfo.ViewNum,
				"labelId":    blogInfo.LabelID,
				"createTime": blog.CreateTime,
				"label":      blog.Label,
			}
			global.RedisSentinel.HMSet(global.CONTEXT, "article:"+blog.Id, mapData)
		} else {
			// 从redis中获取数据
			blog.Id = z.Member.(string)
			fields := []string{"blogId", "title", "describe", "createTime", "views", "label"}
			result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+blog.Id, fields...).Result()

			blog.Id = result[0].(string)
			blog.Title = result[1].(string)
			blog.Content = result[2].(string)
			blog.CreateTime = result[3].(string)
			views := result[4].(string)
			blog.Views, _ = strconv.Atoi(views)
			blog.Label = result[5].(string)
		}
		blogList = append(blogList, blog)
	}

	//生成session_id
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		str := uuid.New().String()
		sessionId = strings.Replace(str, "-", "", -1)

		session := sessions.Default(c)
		option := sessions.Options{MaxAge: 86400, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
		session.Options(option)
		session.Set("session_id", sessionId)
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
	}

	totalCount, err := global.RedisSentinel.ZCard(global.CONTEXT, "blogIdList").Result()
	if err != nil {
		global.LOGGER.Error("无法获取有序集合的总元素数:", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":    totalCount,
		"blogList": blogList,
	})
}

// ReturnArticle 返回具体的文章内容
func ReturnArticle(c *gin.Context) {
	//获取sessionId
	sessionIdString, err := c.Cookie("session_id")
	if err != nil {
		// 没有session_id 则生成session_id
		str := uuid.New().String()
		sessionIdString = strings.Replace(str, "-", "", -1)

		session := sessions.Default(c)
		option := sessions.Options{MaxAge: 86400, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
		session.Options(option)
		session.Set("session_id", sessionIdString)
		if err = session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Msg": "服务器内部错误"})
			return
		}
	}

	// 参数绑定
	var blogId BlogRequest
	if err := c.ShouldBindQuery(&blogId); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"Msg": "参数错误"})
		return
	}

	// 数据获取
	var blog Article
	exist, _ := global.RedisSentinel.Exists(global.CONTEXT, "article:"+blogId.Id).Result()
	if exist == 0 {
		var blogInfo *model.Article
		var effective bool
		// 从mysql中获取文章数据
		id, _ := strconv.ParseUint(blogId.Id, 10, 64)
		if effective, blogInfo = mysql.GetBlogInfoById(id); effective == true {
			blog.Id = blogId.Id
			blog.Title = blogInfo.BlogTitle
			blog.CreateTime = blogInfo.CreatedTime.Format("2006-01-02")
			blog.Views = int(blogInfo.ViewNum)
			blog.Content = blogInfo.Context
			label := mysql.GetLabelNameById(blogInfo.LabelID)
			blog.Label = label
		}

		// 把数据缓存到redis
		mapData := map[string]interface{}{
			"blogId":     blogInfo.ArticleID,
			"title":      blogInfo.BlogTitle,
			"context":    blogInfo.Context,
			"describe":   blogInfo.ShowContext,
			"views":      blogInfo.ViewNum,
			"labelId":    blogInfo.LabelID,
			"createTime": blog.CreateTime,
			"label":      blog.Label,
		}
		global.RedisSentinel.HMSet(global.CONTEXT, "article:"+blog.Id, mapData)
	} else {
		// 从redis中获取数据
		blog.Id = blogId.Id
		fields := []string{"blogId", "title", "context", "createTime", "views", "label"}
		result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+blog.Id, fields...).Result()

		blog.Id = result[0].(string)
		blog.Title = result[1].(string)
		blog.Content = result[2].(string)
		blog.CreateTime = result[3].(string)
		views := result[4].(string)
		blog.Views, _ = strconv.Atoi(views)
		blog.Label = result[5].(string)
	}

	// 浏览量加1
	lua := redis.NewScript(`
	local blogId = KEYS[1]
	local sessionId = KEYS[2]
	local expireTime = tonumber(ARGV[1])

	local isSet = redis.call('SETNX', 'blog:' .. blogId .. ':session:' .. sessionId, 1)
	if isSet == 1 then
		redis.call('EXPIRE', 'blog:' .. blogId .. ':session:' .. sessionId, expireTime)
		redis.call('HINCRBY', 'article:' .. blogId, 'views', 1)
	end
	
	return isSet
	`)

	expireTime := 30
	sessionId := sessionIdString[len(sessionIdString)-19:]
	result, err := lua.Run(global.CONTEXT, global.RedisSentinel, []string{blogId.Id, sessionId}, expireTime).Int()
	if err != nil {
		global.LOGGER.Error("Redis脚本执行错误:", zap.Error(err), zap.Int("isSet的值:", result))
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"article": blog,
	})
}

func ReturnBlogOnceYear(c *gin.Context) {
	// 获取博客id列表
	blogIdList, err := global.RedisSentinel.ZRevRangeWithScores(global.CONTEXT, "blogIdList", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}

	// 用于存储文章的年份和对应的文章列表
	yearGroups := make(map[string][]TimeLineArticle)
	var blog TimeLineArticle

	for _, z := range blogIdList {
		memberStr, ok := z.Member.(string)
		if !ok {
			continue
		}
		fields := []string{"blogId", "title"}
		result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+memberStr, fields...).Result()
		timeInt64 := int64(z.Score)
		createTime := time.Unix(timeInt64, 0).Truncate(time.Second).Format("2006-01-02 15:04:05")
		blogID, ok := result[0].(string)
		if !ok {
			continue
		}
		title, ok := result[1].(string)
		if !ok {
			continue
		}
		blog = TimeLineArticle{
			ID:         blogID,
			Title:      title,
			CreateTime: createTime,
		}
		// 获取文章的创建时间，然后从中解析出年份
		year := strings.Split(blog.CreateTime, "-")[0]
		// 将文章添加到对应年份的文章列表中
		yearGroups[year] = append(yearGroups[year], blog)
	}

	var result []YearGroup
	for year, articles := range yearGroups {
		result = append(result, YearGroup{
			Timestamp: year,
			Articles:  articles,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"blogList": result,
	})
}

func Search(c *gin.Context) {
	//参数绑定
	var questionRequest SearchRequest
	if err := c.ShouldBindQuery(&questionRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	//将输入进行转义（防止XSS攻击）
	question := html.EscapeString(questionRequest.Question)
	//从redis里面查找相关信息
	var title string
	var blog BlogCardCategory
	var blogList []BlogCardCategory

	// 获取博客id列表
	blogIdList, err := global.RedisSentinel.ZRangeWithScores(global.CONTEXT, "blogIdList", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}

	for _, z := range blogIdList {
		if title, err = global.RedisSentinel.HGet(global.CONTEXT, "article:"+z.Member.(string), "title").Result(); err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			} else {
				global.LOGGER.Error("获取文章标题错误:", zap.Error(err))
				continue
			}
		}

		processedTitle := strings.ReplaceAll(title, " ", "")
		processedQuestion := strings.ReplaceAll(question, " ", "")
		processedTitle = strings.ToLower(processedTitle)
		processedQuestion = strings.ToLower(processedQuestion)

		if strings.Contains(processedTitle, processedQuestion) {
			fields := []string{"blogId", "title", "describe", "createTime", "views"}
			result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+z.Member.(string), fields...).Result()

			blog.Id = result[0].(string)
			blog.Title = result[1].(string)
			blog.Content = result[2].(string)
			blog.CreateTime = result[3].(string)
			views := result[4].(string)
			blog.Views, _ = strconv.Atoi(views)

			blogList = append(blogList, blog)
		}
	}
	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"msg":      "关于 '" + question + "' 的搜索结果",
		"blogList": blogList,
	})
}

func GetHotTops(c *gin.Context) {
	var blog HotTops
	var blogList []HotTops
	// 获取博客id列表
	blogIdList, err := global.RedisSentinel.ZRangeWithScores(global.CONTEXT, "blogIdList", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}
	//根据每一个id获取每一篇博客文章的关键信息
	for _, z := range blogIdList {
		fields := []string{"blogId", "title", "views"}
		result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+fmt.Sprint(z.Member), fields...).Result()
		blogId, _ := result[0].(string)
		title, _ := result[1].(string)
		views, _ := result[2].(string)
		blog = HotTops{
			Id:    blogId,
			Title: title,
			Views: views,
		}
		blogList = append(blogList, blog)
	}
	sort.Slice(blogList, func(i, j int) bool {
		viewsI, _ := strconv.Atoi(blogList[i].Views)
		viewsJ, _ := strconv.Atoi(blogList[j].Views)
		return viewsI > viewsJ
	})
	if len(blogList) > 3 {
		blogList = blogList[:3]
	}
	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"hotTops": blogList,
	})
}

func ReturnLabelList(c *gin.Context) {
	//根据有序集合遍历
	var labelInfo Label
	var labelList []Label
	labelSet, err := global.RedisSentinel.ZRevRangeWithScores(global.CONTEXT, "labelSet", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}
	for _, z := range labelSet {
		labelInfo.Name = z.Member.(string)
		labelInfo.ArticleCount = int(z.Score)
		labelList = append(labelList, labelInfo)
	}

	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"labelList": labelList,
	})
}

func Category(c *gin.Context) {
	//参数绑定
	var labelNameRequest LabelName
	if err := c.ShouldBindQuery(&labelNameRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	// 获取博客id列表
	blogIdList, err := global.RedisSentinel.ZRevRangeWithScores(global.CONTEXT, "blogIdList", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}

	var blog BlogCardCategory
	var blogList []BlogCardCategory
	for _, z := range blogIdList {
		labelName, err := global.RedisSentinel.HGet(global.CONTEXT, "article:"+z.Member.(string), "label").Result()
		if errors.Is(err, redis.Nil) {
			continue
		}
		//如果满足匹配条件就放入切片中
		if labelName == labelNameRequest.Label {
			fields := []string{"blogId", "title", "describe", "createTime", "views"}
			result, _ := global.RedisSentinel.HMGet(global.CONTEXT, "article:"+z.Member.(string), fields...).Result()

			blog.Id = result[0].(string)
			blog.Title = result[1].(string)
			blog.Content = result[2].(string)
			blog.CreateTime = result[3].(string)
			views := result[4].(string)
			blog.Views, _ = strconv.Atoi(views)

			blogList = append(blogList, blog)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"blogList": blogList,
	})
}
