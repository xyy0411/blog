package matching

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	"github.com/xyy0411/blog/resp"
)

// GetAllMatchingRecords 获取全部匹配记录
func GetAllMatchingRecords(ctx *gin.Context) {
	var records []models.MatchingRecord
	if err := global.DB.Order("created_at DESC").Find(&records).Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配记录失败")
		return
	}

	resp.OK(ctx, "", map[string]any{
		"total":   len(records),
		"records": records,
	})
}

// GetTodayMatchingRecords 获取今日匹配记录
func GetTodayMatchingRecords(ctx *gin.Context) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var records []models.MatchingRecord
	if err := global.DB.
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询今日匹配记录失败")
		return
	}

	resp.OK(ctx, "", map[string]any{
		"total":   len(records),
		"records": records,
	})
}

// GetThisWeekMatchingRecords 获取本周匹配记录（周一 00:00:00 到当前周最后一天 23:59:59 前）
func GetThisWeekMatchingRecords(ctx *gin.Context) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -(weekday - 1))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var records []models.MatchingRecord
	if err := global.DB.
		Where("created_at >= ? AND created_at < ?", startOfWeek, endOfWeek).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询本周匹配记录失败")
		return
	}

	resp.OK(ctx, "", map[string]any{
		"total":   len(records),
		"records": records,
	})
}
