package models

type Comment struct {
	Model
	UserID    uint   `json:"user_id"`
	Content   string `json:"content"`
	ArticleID uint   `json:"article_id"`
	IP        string `json:"ip"`
}
