package initialize

import (
	"blog/global"
	"fmt"
	"github.com/sony/sonyflake"
	"math/rand"
	"time"
)

func Global() {
	// 初始化 ID 生成器
	rand.Seed(time.Now().Unix())
	startTime, err := time.Parse("2006-01-02 15:04:05", global.StartTime)
	if err != nil {
		fmt.Println(err.Error())
	}
	global.IdGenerator = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: startTime,
	})
}
