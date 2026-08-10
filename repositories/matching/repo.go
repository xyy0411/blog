package matching

import (
	"errors"
	"time"

	"github.com/xyy0411/blog/models"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) UpdateName(userID int64, name string) error {
	return r.db.
		Model(&models.Matching{}).
		Where("user_id = ?", userID).
		Update("name", name).
		Error
}

func (r *Repo) CreateMatchingWithChildren(m *models.Matching) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		if len(m.BlockUsers) > 0 {
			for i := range m.BlockUsers {
				m.BlockUsers[i].MatchingID = m.ID
			}
			if err := tx.Create(&m.BlockUsers).Error; err != nil {
				return err
			}
		}

		if len(m.OnlineSoftwares) > 0 {
			for i := range m.OnlineSoftwares {
				m.OnlineSoftwares[i].MatchingID = m.ID
			}
			return tx.Create(&m.OnlineSoftwares).Error
		}

		return nil
	})
}

func (r *Repo) GetByUserID(userID int64) (models.Matching, error) {
	var m models.Matching
	err := r.db.
		Preload("BlockUsers").
		Preload("OnlineSoftwares").
		Where("user_id = ?", userID).
		First(&m).Error
	return m, err
}

func (r *Repo) UpdateExpire(userID int64, expireAt int64) error {
	return r.db.
		Model(&models.Matching{}).
		Where("user_id = ?", userID).
		Update("expire_at", expireAt).
		Error
}
func (r *Repo) AddOnlineSoftware(matchingID int64, software models.OnlineSoftware) error {
	software.MatchingID = matchingID
	return r.db.Create(&software).Error
}

func (r *Repo) UpdateOnlineSoftwareType(matchingID int64, appName string, appType int8) error {
	return r.db.
		Model(&models.OnlineSoftware{}).
		Where("matching_id = ? AND name = ?", matchingID, appName).
		Update("type", appType).
		Error
}

func (r *Repo) RemoveOnlineSoftware(matchingID int64, appName string) error {
	return r.db.Where("matching_id = ? AND name = ?", matchingID, appName).Delete(&models.OnlineSoftware{}).Error
}

func (r *Repo) RemoveBlockUser(matchingID, userID int64) error {
	return r.db.
		Where("matching_id = ? AND user_id = ?", matchingID, userID).
		Delete(&models.BlockUser{}).
		Error
}

func (r *Repo) AddBlockUser(matchingID, userID int64) error {
	return r.db.Save(&models.BlockUser{
		MatchingID: matchingID,
		UserID:     userID,
	}).Error
}

// ErrInvalidDateRange 表示传入的日期范围无效（开始时间为空或晚于截至时间）
var ErrInvalidDateRange = errors.New("invalid date range: start must not be after end")

// GetApplicationsByDateRange 查询指定日期范围内的 MatchingApplication 记录。
// startDate 为闭区间起点，endDate 为开区间终点，均按本地时间解释。
// 当 startDate 晚于 endDate 时返回 ErrInvalidDateRange 以触发上层参数校验。
func (r *Repo) GetApplicationsByDateRange(startDate, endDate time.Time) ([]models.MatchingApplication, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, ErrInvalidDateRange
	}
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	var apps []models.MatchingApplication
	err := r.db.
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Order("created_at DESC").
		Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}
