package models

import "gorm.io/gorm"

// MatchingRecord 匹配记录表，用于统计用户之间的匹配历史
type MatchingRecord struct {
	gorm.Model
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
