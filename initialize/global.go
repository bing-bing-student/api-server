package initialize

import (
	"blog/global"
	"blog/utils"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func Global() {
	// 初始化 ID 生成器
	rand.New(rand.NewSource(time.Now().Unix()))
	startTime, err := time.Parse("2006-01-02 15:04:05", global.StartTime)
	if err != nil {
		global.LOGGER.Error("Parse time error", zap.Error(err))
	}
	global.IdGenerator = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: startTime,
	})

	//初始化全局翻译器
	if err = utils.Translate("zh"); err != nil {
		global.LOGGER.Error("Translate error", zap.Error(err))
	}
}
