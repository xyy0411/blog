package main

import (
	"github.com/xyy0411/blog/server/config"
	"github.com/xyy0411/blog/server/flag"
	"github.com/xyy0411/blog/server/global"
	"github.com/xyy0411/blog/server/router"
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
