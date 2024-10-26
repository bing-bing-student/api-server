package initialize

import (
	"blog/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// Viper 配置文件的读取
func Viper() {
	// 设置配置文件类型和路径
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/config.yml")

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("获取配置文件错误")
	}

	// 将读取到的配置信息反序列化到全局 CONFIG 中
	if err := viper.Unmarshal(&global.CONFIG); err != nil {
		log.Panic("viper反序列化错误")
	}

	// 监视配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件被修改: ", e.Name)
	})
}
