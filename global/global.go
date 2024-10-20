package global

import (
	"blog/config"
	"context"
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	CONFIG        config.System           // 系统配置信息
	LOGGER        *zap.Logger             // 日志记录
	DB            *gorm.DB                // 数据库接口
	RedisSentinel *redis.Client           // Redis接口（哨兵模式）
	CONTEXT       = context.Background()  // 上下文信息
	IdGenerator   *sonyflake.Sonyflake    // 主键生成器
	Translate     ut.Translator           // 全局翻译器
	StartTime     = "2023-01-01 00:00:01" // 固定启动时间，保证生成 ID 唯一性
	MaxFileSize   = int64(64 << 10)       // MD文件大小限制为64KB
)
