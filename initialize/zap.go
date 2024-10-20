package initialize

import (
	"blog/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func Logger() {
	var level zapcore.Level
	levelStr := global.CONFIG.LogConfig.Level

	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		log.Panicln(err.Error())
	}
	encoder := getEncoder()

	// 创建全量日志的文件切割Writer
	fullLogWriter := getLogWriter(global.CONFIG.LogConfig.Filename)
	core1 := zapcore.NewCore(encoder, fullLogWriter, level)

	// 创建错误日志的文件切割Writer
	errLogWriter := getLogWriter(global.CONFIG.LogConfig.ErrFilename)
	core2 := zapcore.NewCore(encoder, errLogWriter, zap.ErrorLevel)

	// 使用NewTee将core1和core2合并到core
	core := zapcore.NewTee(core1, core2)
	global.LOGGER = zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	// 得到编码配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 通过配置修改时间编码规则
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 通过配置添加调用者信息
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 使用控制台格式保存日志（但是在生产环境下最好是JSON格式）
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
