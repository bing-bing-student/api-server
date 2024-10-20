package redis

import (
	"blog/global"
	"blog/service/mysql"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func WriteViewsOnMysql() {
	// 拿到所有文章id
	blogIdList, err := global.RedisSentinel.ZRevRangeWithScores(global.CONTEXT, "blogIdList", 0, -1).Result()
	if err != nil {
		return
	}
	for _, z := range blogIdList {
		v, _ := global.RedisSentinel.HGet(global.CONTEXT, "article:"+z.Member.(string), "views").Result()
		views, _ := strconv.ParseUint(v, 10, 64)
		id, _ := strconv.ParseUint(z.Member.(string), 10, 64)
		mysql.UpdateViews(id, views)
	}
}

func AddBlogIdToRedis(articleId string, pubTime int64) (err error) {
	ZKey := "blogIdList"
	blogId := []redis.Z{
		{Score: float64(pubTime), Member: articleId},
	}
	err = global.RedisSentinel.ZAdd(global.CONTEXT, ZKey, blogId...).Err()
	return err
}

func AddBlogLabelNameToRedis(labelId uint64) (err error) {
	labelName := mysql.GetLabelNameById(labelId)
	err = global.RedisSentinel.ZIncrBy(global.CONTEXT, "labelSet", 1, labelName).Err()
	return err
}
