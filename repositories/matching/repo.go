package matching

import (
	"github.com/xyy0411/blog/models"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
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
