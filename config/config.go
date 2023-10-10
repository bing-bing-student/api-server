package config

import "time"

//mapstructure用于将通用的map[string]interface{}解码到对应的 Go 结构体中，或者执行相反的操作。

// GinConfig 定义 Gin 配置文件的结构体
type GinConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// RedisConfig 定义 redis 配置文件结构体
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
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

// JWTConfig 定义 jwt 配置文件结构体
type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}

// System 定义项目配置文件结构体
type System struct {
	GinConfig   *GinConfig   `mapstructure:"gin"`
	MySQLConfig *MySQLConfig `mapstructure:"mysql"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
	JWTConfig   *JWTConfig   `mapstructure:"jwt"`
}
