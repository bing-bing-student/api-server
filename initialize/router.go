package initialize

import (
	"blog/controller"
	"blog/global"
	"blog/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()

	//解决跨域问题
	r.Use(cors.Default())

	// 管理员:登录。写文章;改文章;删文章
	admin := r.Group("/admin")
	//admin.POST("/register", controller.Register)
	admin.POST("/login", controller.Login)

	admin.Use(middleware.JWT())
	{
	}

	err := r.Run(fmt.Sprintf("%s:%d", global.CONFIG.GinConfig.Host, global.CONFIG.GinConfig.Port))
	if err != nil {
		fmt.Println(err.Error())
	}
}
