package global

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Logger *log.Logger
)

const (
	// MatchTimeout 匹配超时
	MatchTimeout = 0
	// MatchMsg 匹配中
	MatchMsg = 1
	// MatchSuccess 匹配成功
	MatchSuccess = 2
	// MatchExit 主动退出匹配
	MatchExit = 3
	// MatchError 匹配发生意外错误
	MatchError = 4
	// MatchInteractError 匹配交互时发生错误
	MatchInteractError = 5
)
