package config

import "time"

//mapstructure用于将通用的map[string]interface{}解码到对应的 Go 结构体中，或者执行相反的操作。

// LogConfig 配置日志的结构体
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Filename    string `mapstructure:"filename"`
	ErrFilename string `mapstructure:"err_filename"`
	MaxSize     int    `mapstructure:"max_size"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxBackups  int    `mapstructure:"max_backups"`
}

// GinConfig 定义 Gin 配置文件的结构体
type GinConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// MySQLConfig 定义 mysql 配置文件结构体
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"db_name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig 定义 redis 配置文件结构体
type RedisConfig struct {
	SentinelAdders []string `mapstructure:"sentinel_adders"`
	MasterName     string   `mapstructure:"master_name"`
}

// JWTConfig 定义 jwt 配置文件结构体
type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}

// ALYConfig 阿里云账号配置
type ALYConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
}

// System 定义项目配置文件结构体
type System struct {
	LogConfig   *LogConfig   `mapstructure:"log"`
	GinConfig   *GinConfig   `mapstructure:"gin"`
	MySQLConfig *MySQLConfig `mapstructure:"mysql"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
	JWTConfig   *JWTConfig   `mapstructure:"jwt"`
	ALYConfig   *ALYConfig   `mapstructure:"secret_key"`
}
