package main

import (
	"blog/initialize"
)

func main() {
	// 全局变量初始化
	initialize.Global()

	// 配置信息初始化
	initialize.Viper()

	// 日志初始化
	initialize.Logger()

	// MySQL数据库初始化
	initialize.MySql()

	// Redis缓存初始化
	initialize.Redis()

	// 定时器初始化
	initialize.Timer()

	// 路由初始化
	initialize.Router()
}
