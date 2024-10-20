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

type Blog struct {
	Title    string `form:"title" binding:"required,lte=45"`
	Label    string `form:"label" binding:"required,lte=30"`
	Abstract string `form:"abstract" binding:"required,lte=600"`
}

type BlogArray struct {
	BlogId    string `json:"blog_id"`
	Title     string `json:"title"`
	Label     string `json:"label"`
	Views     uint64 `json:"views"`
	CreatTime string `json:"create_time"`
}

// WriteBlog 编写文章接口
func WriteBlog(c *gin.Context) {
	// md文章长度校验
	bodyLength := c.Request.ContentLength
	if bodyLength > global.MaxFileSize {
		c.JSON(http.StatusOK, gin.H{
			"msg": "文章大小超过64KB",
		})
		return
	}

	type WriteBlogRequest struct {
		Blog Blog   `form:"blog"`
		Text string `form:"text"`
	}

	// 参数绑定
	var blogInfo WriteBlogRequest
	if err := c.ShouldBind(&blogInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	// 校验文章标题是否重复
	title := blogInfo.Blog.Title
	if repeat := mysql.BlogTitleIfExist(title); repeat == true {
		c.JSON(http.StatusOK, gin.H{
			"msg": "文章标题重复了，请换一个标题",
		})
		return
	}
	//拿到标签，看数据库中是否有该标签
	if exits := mysql.LabelIfExist(blogInfo.Blog.Label); exits == 0 {
		labelId, _ := global.IdGenerator.NextID()
		label := model.Label{
			LabelID:   labelId,
			LabelName: blogInfo.Blog.Label,
		}
		if result := mysql.CreateLabel(&label); result == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "提交失败",
			})
			return
		}
	}
	// 插入数据
	local, _ := time.LoadLocation("Asia/Shanghai")
	var article model.Article
	article.ArticleID, _ = global.IdGenerator.NextID()
	article.LabelID = mysql.QueryLabelId(blogInfo.Blog.Label)
	article.BlogTitle = blogInfo.Blog.Title
	article.ShowContext = blogInfo.Blog.Abstract
	article.Context = blogInfo.Text
	article.ViewNum = 0
	article.CreatedTime = time.Now().In(local)
	article.UpdatedTime = article.CreatedTime

	result := mysql.CreateBlog(&article)
	articleId := strconv.FormatUint(article.ArticleID, 10)
	err := redis.AddBlogIdToRedis(articleId, article.CreatedTime.Unix())
	err = redis.AddBlogLabelNameToRedis(article.LabelID)
	if result == 1 && err == nil {

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "提交成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "提交失败",
		})
	}

}

// SortByPubTime 将所有的文章按照发布时间的顺序发送给前端
func SortByPubTime(c *gin.Context) {
	var blogArrayByTime []BlogArray
	for _, v := range mysql.SortBlogByTime() {
		createTime := v.CreatedTime.Format("2006-01-02 15:04:05")
		blogId := strconv.FormatUint(v.ArticleID, 10)
		blog := BlogArray{
			BlogId:    blogId,
			Title:     v.BlogTitle,
			Label:     v.LabelName,
			Views:     v.ViewNum,
			CreatTime: createTime,
		}
		// 将新的 BlogArray 结构体添加到 blogArrayByTime 切片中
		blogArrayByTime = append(blogArrayByTime, blog)
	}
	c.JSON(http.StatusOK, gin.H{
		"blog": blogArrayByTime,
	})
}

// SortByViews 将所有的文章按照浏览量的顺序发送给前端
func SortByViews(c *gin.Context) {
	var blogArrayByViews []BlogArray
	// 这里需要优化
	for _, v := range mysql.SortBlogByViews() {
		blogId := strconv.FormatUint(v.ArticleID, 10)
		createTime := v.CreatedTime.Format("2006-01-02 15:04:05")
		blog := BlogArray{
			BlogId:    blogId,
			Title:     v.BlogTitle,
			Label:     v.LabelName,
			Views:     v.ViewNum,
			CreatTime: createTime,
		}
		// 将新的 BlogArray 结构体添加到 blogArrayByTime 切片中
		blogArrayByViews = append(blogArrayByViews, blog)
	}
	c.JSON(http.StatusOK, gin.H{
		"blog": blogArrayByViews,
	})
}

// GetModifyBlog 根据id返回文章信息
func GetModifyBlog(c *gin.Context) {
	type ModifyBlogRequest struct {
		BlogId string `form:"blog_id" binding:"required"`
	}
	//参数绑定
	var blogId ModifyBlogRequest
	if err := c.ShouldBindQuery(&blogId); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	id, _ := strconv.ParseUint(blogId.BlogId, 10, 64)
	if effective, blogInfo := mysql.GetBlogInfoById(id); effective == true {
		c.JSON(http.StatusOK, gin.H{
			"blog":     blogInfo.Context,
			"describe": blogInfo.ShowContext,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的文章id",
		})
	}
}

// ModifyBlog 修改文章接口
func ModifyBlog(c *gin.Context) {
	// md文章长度校验
	bodyLength := c.Request.ContentLength
	if bodyLength > global.MaxFileSize {
		c.JSON(http.StatusOK, gin.H{
			"msg": "文章大小超过64KB",
		})
		return
	}
	type UpdateBlogRequest struct {
		Blog   Blog   `form:"blog"`
		Text   string `form:"text" binding:"required"`
		BlogId string `json:"blog_id" binding:"required"`
	}

	//参数绑定（文章id,文章标题,文章内容,文章标签,文章描述）
	var blogInfo UpdateBlogRequest
	if err := c.ShouldBind(&blogInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}
	//查看id是否有效
	blogId, _ := strconv.ParseUint(blogInfo.BlogId, 10, 64)
	if effective := mysql.IdIfEffective(blogId); effective == false {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的文章id",
		})
		return
	}

	// 查看是否修改了博客标签
	oldLabelId := mysql.GetLabelIdByBlogId(blogId)
	oldLabelName := mysql.GetLabelNameById(oldLabelId)
	if oldLabelName != blogInfo.Blog.Label {
		// 查看标签是否存在，不存在就在mysql和redis中生成新标签
		if exits := mysql.LabelIfExist(blogInfo.Blog.Label); exits == 0 {
			// 在mysql中新建标签
			labelId, _ := global.IdGenerator.NextID()
			label := model.Label{
				LabelID:   labelId,
				LabelName: blogInfo.Blog.Label,
			}
			if result := mysql.CreateLabel(&label); result == 0 {
				c.JSON(http.StatusOK, gin.H{
					"msg": "更新失败",
				})
				return
			}

			// 在redis当中新建标签,下面的接口已经完成了加1的操作
			if err := redis.AddBlogLabelNameToRedis(labelId); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": "更新失败",
				})
				return
			}
			// 将redis旧标签名的source减一
			err := global.RedisSentinel.ZIncrBy(global.CONTEXT, "labelSet", -1, oldLabelName).Err()
			score, err := global.RedisSentinel.ZScore(global.CONTEXT, "labelSet", oldLabelName).Result()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": "更新失败",
				})
				return
			}
			if score == 0 {
				err = global.RedisSentinel.ZRem(global.CONTEXT, "labelSet", oldLabelName).Err()
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"msg": "更新失败",
					})
					return
				}
			}
		} else {
			// 修改了标签,并且标签存在
			err := global.RedisSentinel.ZIncrBy(global.CONTEXT, "labelSet", 1, blogInfo.Blog.Label).Err()
			err = global.RedisSentinel.ZIncrBy(global.CONTEXT, "labelSet", -1, oldLabelName).Err()
			score, err := global.RedisSentinel.ZScore(global.CONTEXT, "labelSet", oldLabelName).Result()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": "更新失败",
				})
				return
			}
			if score == 0 {
				mysql.DeleteLabel(oldLabelId)
				err = global.RedisSentinel.ZRem(global.CONTEXT, "labelSet", oldLabelName).Err()
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"msg": "更新失败",
					})
					return
				}
			}
		}
	}

	// 更新数据库的数据
	var article model.Article
	// 获取redis当中的文章浏览量
	views, err := global.RedisSentinel.HGet(global.CONTEXT, "article:"+blogInfo.BlogId, "views").Uint64()
	if err != nil {
		// 获取失败，则从mysql中获取
		if views = mysql.GetBlogViewsFromMysql(blogId); views == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "更新失败",
			})
			return
		}
	}

	local, _ := time.LoadLocation("Asia/Shanghai")
	article.UpdatedTime = time.Now().In(local)
	article.LabelID = mysql.QueryLabelId(blogInfo.Blog.Label)
	article.ArticleID = blogId
	article.Context = blogInfo.Text
	article.BlogTitle = blogInfo.Blog.Title
	article.ShowContext = blogInfo.Blog.Abstract
	article.ViewNum = views

	// 延迟双删
	err = global.RedisSentinel.Del(global.CONTEXT, "article:"+blogInfo.BlogId).Err()
	result := mysql.UpdateBlog(&article)
	time.Sleep(300 * time.Millisecond)
	err = global.RedisSentinel.Del(global.CONTEXT, "article:"+blogInfo.BlogId).Err()

	if result != 0 && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新成功",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新失败",
		})
		return
	}
}

// DeleteBlog 删除博客文章
func DeleteBlog(c *gin.Context) {
	type DeleteBlogRequest struct {
		BlogId string `form:"blog_id" binding:"required"`
	}
	//参数绑定
	var blogInfo DeleteBlogRequest
	if err := c.ShouldBind(&blogInfo); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	//查看id是否有效
	blogId, _ := strconv.ParseUint(blogInfo.BlogId, 10, 64)
	if effective := mysql.IdIfEffective(blogId); effective == false {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的文章id",
		})
	}

	labelId := mysql.GetLabelIdByBlogId(blogId)
	labelName := mysql.GetLabelNameById(labelId)

	// blogIdList是一个有序集合,source为创建文章的时间戳,member为博客id
	// labelSet是一个有序集合,source为某一个文章类别下文章的总数量,member为对应的文章标签名称
	err := global.RedisSentinel.ZRem(global.CONTEXT, "blogIdList", blogId).Err()
	err = global.RedisSentinel.ZIncrBy(global.CONTEXT, "labelSet", -1, labelName).Err()
	score, err := global.RedisSentinel.ZScore(global.CONTEXT, "labelSet", labelName).Result()
	if score == 0 {
		err = global.RedisSentinel.ZRem(global.CONTEXT, "labelSet", labelName).Err()
	}

	// 延迟双删
	err = global.RedisSentinel.Del(global.CONTEXT, "article:"+blogInfo.BlogId).Err()
	mysqlSuccess := mysql.DeleteBlog(blogId)
	time.Sleep(300 * time.Millisecond)
	err = global.RedisSentinel.Del(global.CONTEXT, "article:"+blogInfo.BlogId).Err()

	if mysqlSuccess == true && err == nil {
		//检查该文章的标签是否还有其他文章，如果没有其他文章就删除该标签
		mysql.DeleteLabel(labelId)
		c.JSON(http.StatusOK, gin.H{
			"code": "40000",
			"msg":  "删除成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除失败",
		})
	}
}

// GetArticleArrayByLabelId 根据标签id返回该标签的所有博客文章
func GetArticleArrayByLabelId(c *gin.Context) {
	type QueryBlogByLabelRequest struct {
		LabelId string `form:"label_id" binding:"required"`
	}
	//参数绑定
	var labelIdRequest QueryBlogByLabelRequest
	if err := c.ShouldBindQuery(&labelIdRequest); err != nil {
		global.LOGGER.Error("参数绑定错误:", zap.Error(err))
		return
	}

	labelId, _ := strconv.ParseUint(labelIdRequest.LabelId, 10, 64)
	//查询数据
	articleArray := mysql.QueryBlogByLabelId(labelId)

	//返回结果
	type ArticleArrayResponse struct {
		BlogId     string `json:"blog_id"`
		Title      string `json:"title"`
		Views      uint64 `json:"views"`
		CreateTime string `json:"create_time"`
	}
	var articleArrayResponse []ArticleArrayResponse
	for _, article := range articleArray {
		// 创建一个新的ArticleArrayResponse对象，并从Article对象中获取数据
		response := ArticleArrayResponse{
			BlogId:     strconv.FormatUint(article.ArticleID, 10),
			Title:      article.BlogTitle,
			Views:      article.ViewNum,
			CreateTime: article.CreatedTime.Format("2006-01-02 15:04:05"),
		}
		// 将新的ArticleArrayResponse对象添加到切片中
		articleArrayResponse = append(articleArrayResponse, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"articleArray": articleArrayResponse,
	})
}
