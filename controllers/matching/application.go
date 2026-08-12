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

// applicationDateLayout 用于解析前端传入的日期参数，统一为 2006-01-02 形式
const applicationDateLayout = "2006-01-02"

// summarizeApplications 将 MatchingApplication 列表聚合成统计信息，
// 包含总数、成功数、失败数、按退出原因分类统计以及原始记录列表
func summarizeApplications(apps []models.MatchingApplication) map[string]any {
	total := len(apps)
	successCount := 0
	failCount := 0

	// 按 ExitReason 分类统计
	exitReasonStats := map[string]int{
		"success":        0, // ExitReasonSuccess
		"user_initiated": 0, // ExitReasonUserInitiated
		"timeout":        0, // ExitReasonTimeout
		"error":          0, // ExitReasonError
		"expired":        0, // ExitReasonExpired
	}

	for _, app := range apps {
		if app.IsMatched {
			successCount++
		} else {
			failCount++
		}

		// 统计各个 ExitReason
		switch app.ExitReason {
		case models.ExitReasonSuccess:
			exitReasonStats["success"]++
		case models.ExitReasonUserInitiated:
			exitReasonStats["user_initiated"]++
		case models.ExitReasonTimeout:
			exitReasonStats["timeout"]++
		case models.ExitReasonError:
			exitReasonStats["error"]++
		case models.ExitReasonExpired:
			exitReasonStats["expired"]++
		}
	}

	return map[string]any{
		"total":         total,
		"success_count": successCount,
		"fail_count":    failCount,
		"success_rate":  computeRate(successCount, total),
		"fail_rate":     computeRate(failCount, total),
		"exit_reason_stats": map[string]any{
			"success":        map[string]any{"count": exitReasonStats["success"], "rate": computeRate(exitReasonStats["success"], total)},
			"user_initiated": map[string]any{"count": exitReasonStats["user_initiated"], "rate": computeRate(exitReasonStats["user_initiated"], total)},
			"timeout":        map[string]any{"count": exitReasonStats["timeout"], "rate": computeRate(exitReasonStats["timeout"], total)},
			"error":          map[string]any{"count": exitReasonStats["error"], "rate": computeRate(exitReasonStats["error"], total)},
			"expired":        map[string]any{"count": exitReasonStats["expired"], "rate": computeRate(exitReasonStats["expired"], total)},
		},
		"records": apps,
	}
}

// computeRate 计算百分比，避免除零；保留两位小数对应的浮点值
func computeRate(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

// GetTodayMatchingApplications 查询当日 MatchingApplication 表数据，
// 无需日期参数，默认使用系统当前日期
func GetTodayMatchingApplications(ctx *gin.Context) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	repo := matchingrepo.NewRepo(global.DB)
	apps, err := repo.GetApplicationsByDateRange(startOfDay, endOfDay)
	if err != nil {
		global.Logger.Errorf("查询当日 MatchingApplication 失败: %v", err)
		resp.Error(ctx, http.StatusInternalServerError, "查询当日匹配申请数据失败")
		return
	}

	resp.OK(ctx, "", summarizeApplications(apps))
}

// GetMatchingApplicationsByDateRange 根据 start_date 与 end_date 查询指定日期范围内的数据
func GetMatchingApplicationsByDateRange(ctx *gin.Context) {
	startStr := ctx.Query("start_date")
	endStr := ctx.Query("end_date")

	if startStr == "" || endStr == "" {
		resp.Error(ctx, http.StatusBadRequest, "start_date 和 end_date 参数不能为空")
		return
	}

	loc := time.Now().Location()
	startDate, err := time.ParseInLocation(applicationDateLayout, startStr, loc)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "start_date 格式无效，应为 YYYY-MM-DD")
		return
	}

	endDate, err := time.ParseInLocation(applicationDateLayout, endStr, loc)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "end_date 格式无效，应为 YYYY-MM-DD")
		return
	}

	// 将 endDate 推进到次日 00:00:00，使其成为闭区间的右端点
	endDate = endDate.Add(24 * time.Hour)

	repo := matchingrepo.NewRepo(global.DB)
	apps, err := repo.GetApplicationsByDateRange(startDate, endDate)
	if err != nil {
		if errors.Is(err, matchingrepo.ErrInvalidDateRange) {
			resp.Error(ctx, http.StatusBadRequest, "开始日期不能晚于截至日期")
			return
		}
		global.Logger.Errorf("查询 MatchingApplication 日期范围失败: %v", err)
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配申请数据失败")
		return
	}

	resp.OK(ctx, "", summarizeApplications(apps))
}
