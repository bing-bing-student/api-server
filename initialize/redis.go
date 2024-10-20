package initialize

import (
	"blog/global"
	"github.com/redis/go-redis/v9"
)

func Redis() {
	// 连接redis哨兵集群
	masterName := global.CONFIG.RedisConfig.MasterName
	sentinelAdders := global.CONFIG.RedisConfig.SentinelAdders

	// 默认情况下连接池的大小是:runtime.GOMAXPROCS * 10 。本机中连接池大小就是80
	global.RedisSentinel = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAdders,
	})

	// 检查 RedisSentinel 连通性
	if _, err := global.RedisSentinel.Ping(global.CONTEXT).Result(); err != nil {
		panic(err.Error())
	}
}
