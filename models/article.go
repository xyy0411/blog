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
	OpenComment  bool   `json:"open_comment"`         // 文章评论开关
	// Status       int8   `json:"status"`               // 状态 草稿 审核中 已发布
}

// CommitArticle 前端向后端提交时需要的数据
type CommitArticle struct {
	Title       string         `gorm:"size:32" json:"title"` // 标题
	Abstract    string         `json:"abstract"`             // 简介
	Cover       string         `json:"cover"`                // 封面
	OpenComment bool           `json:"open_comment"`         // 文章评论开关
	Contents    ArticleContent `json:"contents"`             // 文章内容 因为要分段所以是是切片
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

type ArticleContent struct {
	Content []string `json:"content"`
}

func (c ArticleContent) Len() int {
	return len(c.Content)
}

func (c ArticleContent) String() string {
	var str string
	for _, v := range c.Content {
		str += v + "\n"
	}
	return str
}
