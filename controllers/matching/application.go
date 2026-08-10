package matching

import (
	"errors"
	"fmt"
	"math/rand"
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
// 包含总数、成功数、失败数以及原始记录列表
func summarizeApplications(apps []models.MatchingApplication) map[string]any {
	total := len(apps)
	successCount := 0
	failCount := 0
	for _, app := range apps {
		if app.IsMatched {
			successCount++
		} else {
			failCount++
		}
	}

	return map[string]any{
		"total":         total,
		"success_count": successCount,
		"fail_count":    failCount,
		"success_rate":  computeRate(successCount, total),
		"fail_rate":     computeRate(failCount, total),
		"records":       apps,
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

// SeedMatchingApplications 向 matching_applications 表中插入一批模拟数据。
// 仅用于本地演示/验证效果：随机挑选近 10 天，生成约 120 条数据（约 60% 成功，40% 失败）。
// 接口幂等：每次调用都会追加新数据，不会覆盖已有。
func SeedMatchingApplications(ctx *gin.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	users := []struct {
		id   int64
		name string
	}{
		{1001, "小明"},
		{1002, "小红"},
		{1003, "张三"},
		{1004, "李四"},
		{1005, "王五"},
		{1006, "赵六"},
		{1007, "阿七"},
		{1008, "小八"},
	}

	now := time.Now()
	loc := now.Location()
	createdCount := 0
	successCount := 0
	failCount := 0

	// 每天按权重生成条数：越靠近今天条数越多
	for daysAgo := 9; daysAgo >= 0; daysAgo-- {
		dayBase := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
			AddDate(0, 0, -daysAgo)

		// 近 10 天的总目标约 120 条，前几天少，今天多
		weight := 1.0
		switch {
		case daysAgo == 0:
			weight = 2.4 // 今天 24 条左右
		case daysAgo <= 2:
			weight = 1.6
		case daysAgo <= 5:
			weight = 1.1
		default:
			weight = 0.7
		}
		perDayTarget := int(float64(10) * weight)

		for i := 0; i < perDayTarget; i++ {
			user := users[rng.Intn(len(users))]

			// 62% 的成功率
			isMatched := rng.Intn(100) < 62

			var duration int
			var matchID string
			if isMatched {
				duration = 30 + rng.Intn(270) // 30~300 秒
				matchID = fmt.Sprintf("M%d%d%d", dayBase.Unix(), user.id, rng.Intn(9000)+1000)
				successCount++
			} else {
				duration = 0
				failCount++
			}

			// 在当天内随机偏移时间（0~23 小时 + 0~59 分钟 + 0~59 秒）
			offset := time.Duration(rng.Intn(24))*time.Hour +
				time.Duration(rng.Intn(60))*time.Minute +
				time.Duration(rng.Intn(60))*time.Second
			createdAt := dayBase.Add(offset)

			app := models.MatchingApplication{
				UserID:    user.id,
				UserName:  user.name,
				IsMatched: isMatched,
				Duration:  duration,
				MatchID:   matchID,
			}
			app.CreatedAt = createdAt
			app.UpdatedAt = createdAt

			if err := global.DB.Create(&app).Error; err != nil {
				global.Logger.Errorf("插入 MatchingApplication 测试数据失败: %v", err)
				resp.Error(ctx, http.StatusInternalServerError, "生成测试数据失败，请检查数据库连接")
				return
			}
			createdCount++
		}
	}

	resp.OK(ctx, "已生成匹配申请测试数据", map[string]any{
		"created_count":  createdCount,
		"success_count":  successCount,
		"fail_count":     failCount,
		"success_rate":   computeRate(successCount, createdCount),
		"date_span":      "最近 10 天（含今日）",
		"distinct_users": len(users),
	})
}
