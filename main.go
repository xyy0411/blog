package main

import (
	"github.com/xyy0411/blog/config"
	"github.com/xyy0411/blog/flag"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/router"
)

func main() {
	flag.Parse()
	config.InitLog()
	config.InitConfig()
	flag.Run()
	engine := router.SetupRouter()
	err := engine.Run(config.AppConfig.App.Port)
	if err != nil {
		global.Logger.Error("启动失败:", err)
		return
	}
}
