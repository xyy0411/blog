package matching

import (
	"strconv"
	"time"

	"github.com/RomiChan/syncx"
	"github.com/xyy0411/blog/config"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	"github.com/xyy0411/blog/utils"
)

const defaultMatchCooldownMinutes = 30

type Manager struct {
	matchedList syncx.Map[int64, *models.Matching]
	idGen       *Snowflake
}

func NewMatchingManager() *Manager {
	return &Manager{
		matchedList: syncx.Map[int64, *models.Matching]{},
		idGen:       NewSnowflake(0),
	}
}

func (mm *Manager) Len() int {
	var listLen int

	mm.matchedList.Range(func(key int64, value *models.Matching) bool {
		listLen++
		return true
	})

	return listLen
}

func (mm *Manager) AddUserToQueue(user *models.Matching) {
	mm.matchedList.Store(user.UserID, user)
}

func (mm *Manager) RemoveUserFromQueue(userID int64) {
	mm.matchedList.Delete(userID)
}

func (mm *Manager) GenerateMatchID() string {
	return strconv.FormatInt(mm.idGen.NextID(), 10)
}

func (mm *Manager) matchCooldownMinutes() int {
	if config.AppConfig == nil || config.AppConfig.Matching.CooldownMinutes <= 0 {
		return defaultMatchCooldownMinutes
	}
	return config.AppConfig.Matching.CooldownMinutes
}

func (mm *Manager) inMatchCooldown(userID, targetUserID int64) bool {
	if global.DB == nil {
		return false
	}

	cutoff := time.Now().Add(-time.Duration(mm.matchCooldownMinutes()) * time.Minute)
	var recentCount int64
	err := global.DB.Model(&models.MatchingRecord{}).
		Where(
			"created_at >= ? AND ((user_id = ? AND peer_id = ?) OR (user_id = ? AND peer_id = ?))",
			cutoff, userID, targetUserID, targetUserID, userID,
		).
		Count(&recentCount).Error
	if err != nil {
		global.Logger.Errorf("查询匹配冷却期失败: %v", err)
		return false
	}

	return recentCount > 0
}

func (mm *Manager) saveMatchingRecord(user, targetUser models.Matching, matchID string) {
	record := models.MatchingRecord{
		UserID:   user.UserID,
		UserName: user.UserName,
		PeerID:   targetUser.UserID,
		PeerName: targetUser.UserName,
		MatchID:  matchID,
	}

	if err := global.DB.Create(&record).Error; err != nil {
		global.Logger.Errorf("保存匹配记录失败: %v", err)
	}
}

// saveMatchingApplication 保存匹配申请记录到数据库
func (mm *Manager) saveMatchingApplication(userID int64, userName string, isMatched bool, duration int, matchID string, exitReason models.ExitReason) {
	app := models.MatchingApplication{
		UserID:     userID,
		UserName:   userName,
		IsMatched:  isMatched,
		Duration:   duration,
		MatchID:    matchID,
		ExitReason: exitReason,
	}

	if err := global.DB.Create(&app).Error; err != nil {
		global.Logger.Errorf("保存匹配申请记录失败: userID=%d, err=%v", userID, err)
	}
}

func (mm *Manager) notifyAndRemoveUser(id int64, user models.Matching, matchID string) {
	if MatchHub == nil {
		global.Logger.Error("MatchHub 未初始化")
		return
	}
	if MatchHub.clients == nil {
		global.Logger.Error("MatchHub.clients 未初始化")
		return
	}
	client, ok := MatchHub.clients[id]
	if !ok || client == nil {
		global.Logger.Errorf("用户 %d 的客户端未找到", user.UserID)
		return
	}
	if client.send == nil {
		global.Logger.Errorf("用户 %d 的客户端 send 通道未初始化", user.UserID)
		return
	}

	// 计算匹配持续时间（秒）
	duration := int(time.Now().Sub(client.connectedAt).Seconds())

	// 保存匹配申请记录（匹配成功）
	mm.saveMatchingApplication(id, user.UserName, true, duration, matchID, models.ExitReasonSuccess)

	event := utils.FormatMatchingInfo(id, user, matchID)
	mm.RemoveUserFromQueue(user.UserID)
	sendEvent(client, event)
	client.close <- true
}

func (mm *Manager) MatchUsers(user models.Matching) {
	if mm.Len() == 0 {
		mm.AddUserToQueue(&user)
		return
	}

	var targetUser models.Matching
	mm.matchedList.Range(func(key int64, value *models.Matching) bool {
		if user.UserID != value.UserID && user.IsMatch(*value) && !mm.inMatchCooldown(user.UserID, value.UserID) {
			targetUser = *value
			return false
		}
		return true
	})

	if targetUser.UserID == 0 {
		mm.AddUserToQueue(&user)
		return
	}

	matchID := mm.GenerateMatchID()
	mm.notifyAndRemoveUser(targetUser.UserID, user, matchID)
	mm.notifyAndRemoveUser(user.UserID, targetUser, matchID)

	global.Logger.Infof("匹配成功 用户:%d <----> 用户:%d, 匹配ID:%s", user.UserID, targetUser.UserID, matchID)
	mm.saveMatchingRecord(user, targetUser, matchID)
}
