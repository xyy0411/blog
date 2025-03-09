package global

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

var (
	DB           *gorm.DB
	Logger       *log.Logger
	RegisterLock sync.RWMutex
)
