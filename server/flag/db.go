package flag

import (
	"github.com/xyy0411/blog/server/global"
	"github.com/xyy0411/blog/server/models"
	"os"
)

func flagDB() {
	err := global.DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.ArticleLikes{},
		&models.Comment{},
		&models.Banner{},
		&models.Log{},
		&models.UserLogin{},
	)
	if err != nil {
		global.Logger.Errorf("数据库迁移失败: %v", err)
		return
	}

	global.Logger.Info("数据库迁移成功!")
}

func Run() {
	if flagOptions.DB {
		flagDB()
		os.Exit(0)
	}
}
