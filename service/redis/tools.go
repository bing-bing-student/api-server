package redis

import (
	"blog/global"
	"strconv"
)

func AddToolsIdToRedis(toolId uint64) (err error) {
	id := strconv.FormatUint(toolId, 10)
	err = global.RedisSentinel.SAdd(global.CONTEXT, "toolsSet", id).Err()
	return err
}
