package mysql

import (
	"blog/global"
	"blog/model"
	"time"
)

// Result 博客文章查询
type Result struct {
	ArticleID   uint64
	BlogTitle   string
	LabelName   string
	CreatedTime time.Time
	ViewNum     uint64
}

func BlogTitleIfExist(title string) bool {
	var blog *model.Article
	result := global.DB.Where("blog_title=?", title).First(&blog)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func CreateBlog(article *model.Article) int64 {
	result := global.DB.Create(&article).RowsAffected
	return result
}

func GetBlogViewsFromMysql(blogId uint64) uint64 {
	var article model.Article
	if err := global.DB.Where("article_id=?", blogId).First(&article).Error; err != nil {
		return 0
	} else {
		return article.ViewNum
	}
}

func UpdateViews(id, views uint64) {
	var article model.Article
	global.DB.Model(&article).Where("article_id=?", id).UpdateColumn("view_num", views)
}

func SortBlogByTime() (resultArray []Result) {
	global.DB.Raw(
		"SELECT A.article_id, A.blog_title, L.label_name, A.created_time, A.view_num " +
			"FROM article AS A " +
			"JOIN label AS L " +
			"ON A.label_id = L.label_id " +
			"ORDER BY A.created_time DESC").Scan(&resultArray)
	return resultArray
}

func SortBlogByViews() (resultArray []Result) {
	global.DB.Raw(
		"SELECT A.article_id, A.blog_title, L.label_name, A.created_time, A.view_num " +
			"FROM article AS A " +
			"JOIN label AS L " +
			"ON A.label_id = L.label_id " +
			"ORDER BY A.view_num DESC").Scan(&resultArray)
	return resultArray
}

func GetBlogInfoById(id uint64) (effective bool, blogInfo *model.Article) {
	result := global.DB.Where("article_id=?", id).First(&blogInfo).RowsAffected
	if result == 0 {
		effective = false
		blogInfo = nil
		return effective, blogInfo
	} else {
		effective = true
		return effective, blogInfo
	}
}

func IdIfEffective(id uint64) (effective bool) {
	var blogInfo *model.Article
	result := global.DB.Where("article_id=?", id).First(&blogInfo).RowsAffected
	if result == 1 {
		effective = true
		return effective
	}
	effective = false
	return effective
}

func UpdateBlog(article *model.Article) (result int64) {
	result = global.DB.Model(&article).Where("article_id=?", article.ArticleID).Updates(article).RowsAffected
	return result
}

func DeleteBlog(id uint64) (success bool) {
	var blog *model.Article
	if result := global.DB.Delete(&blog, "article_id=?", id).RowsAffected; result == 1 {
		success = true
		return success
	}
	success = false
	return success
}

func GetLabelIdByBlogId(id uint64) (labelId uint64) {
	var blog *model.Article
	if result := global.DB.Where("article_id=?", id).First(&blog).RowsAffected; result == 1 {
		labelId = blog.LabelID
		return labelId
	}
	return labelId

}

func QueryBlogByLabelId(id uint64) (blogArray []*model.Article) {
	global.DB.Where("label_id=?", id).Find(&blogArray)
	return blogArray
}
