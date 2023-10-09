package global

import (
	"blog/config"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

var (
	CONFIG      config.System           // 系统配置信息
	DB          *gorm.DB                // 数据库接口
	REDIS       *redis.Client           // Redis 缓存接口
	CONTEXT     = context.Background()  // 上下文信息
	IdGenerator *sonyflake.Sonyflake    // 主键生成器
	StartTime   = "2023-01-01 00:00:01" // 固定启动时间，保证生成 ID 唯一性
	MaxFileSize = int64(64 << 10)       // MD文件大小限制为64KB
)
