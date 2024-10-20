package initialize

import (
	"blog/global"
	"blog/service/redis"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func Timer() {
	cronTab := cron.New()

	// 定时任务,cron表达式
	spec := "0 2,13,19 * * *"
	// 添加定时任务
	_, err := cronTab.AddFunc(spec, redis.WriteViewsOnMysql)
	if err != nil {
		global.LOGGER.Error("定时任务错误:", zap.Error(err))
		return
	}
	// 启动定时器
	cronTab.Start()
}
