package models

import "time"

type Article struct {
	Model
	UserID       uint   `json:"user_id"`              // 关联用户
	Title        string `gorm:"size:32" json:"title"` // 标题
	Abstract     string `json:"abstract"`             // 简介
	Content      string `json:"content"`              // 内容
	Cover        string `json:"cover"`                // 封面
	LookCount    int    `json:"look_count"`           // 浏览量
	Likes        int    `json:"likes"`                // 点赞数
	CommentCount int    `json:"comment_count"`        // 评论数
	CollectCount int    `json:"collect_count"`        // 收藏数
	OpenComment  int    `json:"open_comment"`         // 文章评论开关
	Status       int8   `json:"status"`               // 状态 草稿 审核中 已发布
}

type CommentWithArticle struct {
	Article  Article    `json:"article"`
	Comments []*Comment `json:"comments"`
}

type ArticleLikes struct {
	UserID       uint      `gorm:"uniqueIndex:idx_name" json:"user_id"`
	UserModel    User      `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint      `gorm:"uniqueIndex:idx_name" json:"article_id"`
	ArticleModel Article   `gorm:"foreignKey:ArticleID" json:"-"`
	CreatedAt    time.Time `json:"created_at"` // 点赞时间
}
