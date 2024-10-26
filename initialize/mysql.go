package initialize

import (
	"blog/global"
	"blog/model"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func MySql() {
	username := global.CONFIG.MySQLConfig.Username // 账号
	password := global.CONFIG.MySQLConfig.Password // 密码
	host := global.CONFIG.MySQLConfig.Host         // 数据库地址，可以是Ip或者域名
	port := global.CONFIG.MySQLConfig.Port         // 数据库端口
	dbName := global.CONFIG.MySQLConfig.DBName     // 数据库名
	// dsn := "用户名:密码@tcp(地址:端口)/数据库名"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)

	// 配置Gorm连接到MySQL
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}

	//日志
	fileLogger := &lumberjack.Logger{
		Filename:   "./log/mysql.log",
		MaxSize:    2,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	newLogger := logger.New(
		log.New(fileLogger, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 0, // 慢 SQL 阈值
			LogLevel:                  logger.Info,          // 日志级别
			IgnoreRecordNotFoundError: true,                 // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                // 禁用彩色打印
		},
	)
	//禁用复数表名
	newNamingStrategy := schema.NamingStrategy{
		SingularTable: true,
	}
	if DB, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger, NamingStrategy: newNamingStrategy}); err == nil {
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxOpenConns(global.CONFIG.MySQLConfig.MaxOpenConns)       // 设置数据库最大连接数
		sqlDB.SetMaxIdleConns(global.CONFIG.MySQLConfig.MaxIdleConns)       // 设置上数据库最大闲置连接数
		sqlDB.SetConnMaxLifetime(global.CONFIG.MySQLConfig.ConnMaxLifetime) // 设置连接可复用的最大时间
		global.DB = DB
	} else {
		panic("connect server failed")
	}

	// 自动生成对应的数据库表
	err := global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Article{})
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Tools{})
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Label{})
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.UserInfo{})
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.UserLogin{})
	if err != nil {
		global.LOGGER.Error("表结构错误:", zap.Error(err))
	}
}
