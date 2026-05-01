package matching

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	matchingRepo "github.com/xyy0411/blog/repositories/matching"
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

// GetMatchingApplicationStats 展示匹配申请统计（当日/任意时间/总和）
func GetMatchingApplicationStats(ctx *gin.Context) {
	scope := ctx.DefaultQuery("scope", "all")
	var filter matchingRepo.MatchingApplicationListFilter
	now := time.Now()

	switch scope {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.Add(24 * time.Hour)
		filter.StartAt = &start
		filter.EndAt = &end
	case "range":
		startStr := ctx.Query("start")
		endStr := ctx.Query("end")
		if startStr == "" || endStr == "" {
			resp.Error(ctx, http.StatusBadRequest, "scope=range 时必须提供 start 和 end，格式为 RFC3339")
			return
		}
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			resp.Error(ctx, http.StatusBadRequest, "start 格式错误，需为 RFC3339")
			return
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			resp.Error(ctx, http.StatusBadRequest, "end 格式错误，需为 RFC3339")
			return
		}
		filter.StartAt = &start
		filter.EndAt = &end
	case "all":
	default:
		resp.Error(ctx, http.StatusBadRequest, "scope 仅支持 today/range/all")
		return
	}

	repo := matchingRepo.NewRepo(global.DB)
	total, matched, successRate, durationSum, err := repo.MatchingApplicationStats(filter)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配申请统计失败")
		return
	}
	resp.OK(ctx, "", gin.H{
		"scope":                 scope,
		"total_count":           total,
		"matched_count":         matched,
		"success_rate":          successRate,
		"total_duration_second": durationSum,
	})
}
