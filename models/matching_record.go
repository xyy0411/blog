package models

import (
	"time"

	"gorm.io/gorm"
)

// MatchingRecord 匹配记录表，用于统计用户之间的匹配历史
type MatchingRecord struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UserID   int64  `json:"user_id"`   // 发起匹配的用户ID
	UserName string `json:"user_name"` // 发起匹配的用户名
	PeerID   int64  `json:"peer_id"`   // 被匹配的用户ID
	PeerName string `json:"peer_name"` // 被匹配的用户名
	MatchID  string `json:"match_id"`  // 匹配ID
}

// TableName 指定表名
func (MatchingRecord) TableName() string {
	return "matching_records"
}

// MatchingApplication 匹配申请表，记录用户申请匹配的信息
type MatchingApplication struct {
	gorm.Model
	UserID    int64  `gorm:"index" json:"user_id"`                    // 用户ID
	UserName  string `gorm:"size:100" json:"user_name"`               // 用户名
	IsMatched bool   `gorm:"default:false;index" json:"is_matched"`   // 是否匹配成功
	Duration  int    `gorm:"default:0" json:"duration"`               // 匹配持续时间（秒）
	MatchID   string `gorm:"size:50;index" json:"match_id,omitempty"` // 匹配ID
}

// TableName 指定表名
func (MatchingApplication) TableName() string {
	return "matching_applications"
}
