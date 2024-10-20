package initialize

import (
	"blog/controller"
	"blog/global"
	"blog/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Router() {
	r := gin.New()
	// 禁用代理
	if err := r.SetTrustedProxies(nil); err != nil {
		return
	}

	// 解决跨域问题
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://liubing.xyz", "https://bingbingstudent-0929-admin.liubing.xyz"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}

	// 使用中间件
	r.Use(cors.New(corsConfig), middleware.GinLogger(), middleware.GinRecovery(true))

	// 后台管理
	admin := r.Group("/admin")
	admin.POST("/return-code", controller.ReturnCode)
	admin.POST("/short-message-login", controller.LoginBySMS)
	admin.POST("/refresh-token", controller.RefreshToken)
	admin.Use(middleware.JWT())
	{
		admin.POST("/login", controller.LoginByToken)
		admin.POST("/writeBlog", controller.WriteBlog)
		admin.GET("/sortByPubTime", controller.SortByPubTime)
		admin.GET("/sortByViews", controller.SortByViews)
		admin.GET("/getModifyBlog", controller.GetModifyBlog)
		admin.POST("/modifyBlog", controller.ModifyBlog)
		admin.DELETE("/deleteBlog", controller.DeleteBlog)
		admin.GET("/getAllLabel", controller.GetAllLabel)
		admin.GET("/articleArray", controller.GetArticleArrayByLabelId)
		admin.GET("/getTools", controller.GetTools)
		admin.POST("/addTools", controller.AddTools)
		admin.PUT("/modifyTools", controller.UpdateTools)
		admin.DELETE("/deleteTools", controller.DeleteTools)
		admin.PUT("/aboutMe", controller.ManageAboutMe)
	}

	// 前台展示
	user := r.Group("/user")
	user.Use(sessions.Sessions("session_id", cookie.NewStore([]byte("secret"))))
	{
		user.GET("/index", controller.GetBlogListOnIndex)
		user.GET("/search", controller.Search)
		user.GET("/hotTops", controller.GetHotTops)
		user.GET("/getArticle", controller.ReturnArticle)
		user.GET("/getLabelList", controller.ReturnLabelList)
		user.GET("/category", controller.Category)
		user.GET("/getToolsList", controller.ReturnToolsList)
		user.GET("/getBlogOnceYear", controller.ReturnBlogOnceYear)
		user.GET("/getAboutMeInfo", controller.ReturnAboutMeInfo)
	}

	err := r.Run(fmt.Sprintf("%s:%d", global.CONFIG.GinConfig.Host, global.CONFIG.GinConfig.Port))
	if err != nil {
		global.LOGGER.Error("路由错误:", zap.Error(err))
	}
}
