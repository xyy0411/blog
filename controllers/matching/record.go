package matching

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	matchingrepo "github.com/xyy0411/blog/repositories/matching"
	"github.com/xyy0411/blog/resp"
)

// recordDateLayout 用于解析前端传入的日期参数，统一为 2006-01-02 形式
const recordDateLayout = "2006-01-02"

// getRecordsByDateRange 根据日期范围查询 MatchingRecord
func getRecordsByDateRange(startDate, endDate time.Time) ([]models.MatchingRecord, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, matchingrepo.ErrInvalidDateRange
	}
	if startDate.After(endDate) {
		return nil, matchingrepo.ErrInvalidDateRange
	}

	var records []models.MatchingRecord
	if err := global.DB.
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

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

	records, err := getRecordsByDateRange(startOfDay, endOfDay)
	if err != nil {
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

	records, err := getRecordsByDateRange(startOfWeek, endOfWeek)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询本周匹配记录失败")
		return
	}

	resp.OK(ctx, "", map[string]any{
		"total":   len(records),
		"records": records,
	})
}

// GetMatchingRecordsByDateRange 根据 start_date 与 end_date 查询指定日期范围内的匹配记录
func GetMatchingRecordsByDateRange(ctx *gin.Context) {
	startStr := ctx.Query("start_date")
	endStr := ctx.Query("end_date")

	if startStr == "" || endStr == "" {
		resp.Error(ctx, http.StatusBadRequest, "start_date 和 end_date 参数不能为空")
		return
	}

	loc := time.Now().Location()
	startDate, err := time.ParseInLocation(recordDateLayout, startStr, loc)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "start_date 格式无效，应为 YYYY-MM-DD")
		return
	}

	endDate, err := time.ParseInLocation(recordDateLayout, endStr, loc)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "end_date 格式无效，应为 YYYY-MM-DD")
		return
	}

	// 将 endDate 推进到次日 00:00:00，使其成为闭区间的右端点
	endDate = endDate.Add(24 * time.Hour)

	records, err := getRecordsByDateRange(startDate, endDate)
	if err != nil {
		if errors.Is(err, matchingrepo.ErrInvalidDateRange) {
			resp.Error(ctx, http.StatusBadRequest, "开始日期不能晚于截至日期")
			return
		}
		global.Logger.Errorf("查询 MatchingRecord 日期范围失败: %v", err)
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配记录数据失败")
		return
	}

	resp.OK(ctx, "", map[string]any{
		"total":   len(records),
		"records": records,
	})
}
