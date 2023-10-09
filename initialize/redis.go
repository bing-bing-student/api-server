package initialize

import (
	"blog/global"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Redis() {
	host := global.CONFIG.RedisConfig.Host
	port := global.CONFIG.RedisConfig.Port
	db := global.CONFIG.RedisConfig.DB
	poolSize := global.CONFIG.RedisConfig.PoolSize

	addr := fmt.Sprintf("%s:%d", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // redis地址
		DB:       db,   // redis数据库编号，默认为0
		PoolSize: poolSize,
	})
	// 检查 Redis 连通性
	if _, err := rdb.Ping(global.CONTEXT).Result(); err != nil {
		panic(err.Error())
	}
	global.REDIS = rdb
}
