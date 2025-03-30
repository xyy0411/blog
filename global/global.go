package global

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Logger *log.Logger
)
