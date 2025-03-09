package config

import (
	"github.com/xyy0411/blog/server/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func initDB() {
	db, err := gorm.Open(postgres.Open(AppConfig.Database.Dsn), &gorm.Config{})

	if err != nil {
		global.Logger.Error("数据库连接失败")
		return
	}

	sqlDB, err := db.DB()

	if err != nil {
		global.Logger.Error("数据库连接失败")
		return
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleCons)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenCons)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	global.Logger.Info("数据库连接成功")
}
