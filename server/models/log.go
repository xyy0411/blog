package models

type Log struct {
	Model
	LogType   int8   `json:"log_type"` // 日志类型 1 2 3
	Title     string `json:"title"`
	Content   string `json:"content"`
	Level     int8   `json:"level"` // 日志等级 1 2 3
	UserID    uint   `json:"user_id"`
	UserModel User   `gorm:"foreignKey:UserID" json:"-"` // 用户信息
	IP        string `json:"ip"`
	Address   string `json:"address"`
	IsRead    bool   `json:"is_read"` // 是否读取
}
