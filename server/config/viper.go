package config

import (
	"github.com/spf13/viper"
	"github.com/xyy0411/blog/server/global"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn         string
		MaxOpenCons int
		MaxIdleCons int
	}
}

var AppConfig *Config

func InitConfig() {
	viper.AddConfigPath("./server/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		global.Logger.Error("读取config.yaml文件失败:", err)
		viper.AddConfigPath("./config")
		err = viper.ReadInConfig()
		if err != nil {
			global.Logger.Error(err)
			return
		}
	}
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		global.Logger.Error("解析config.yaml文件错误:", err)
		return
	}
	global.Logger.Info("读取config.yaml文件成功")
	initDB()
}
