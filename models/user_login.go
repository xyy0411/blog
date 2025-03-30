package models

type UserLogin struct {
	Model
	UserID    uint   `json:"user_id"`
	UserModel User   `gorm:"foreignKey:UserID" json:"-"` // 用户信息
	IP        string `json:"ip"`
	Address   string `json:"address"`
	UA        string `json:"ua"`
}
