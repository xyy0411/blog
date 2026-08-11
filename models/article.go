package models

import "time"

type Article struct {
        Model
        UserID       uint   `json:"user_id"`
        Title        string `gorm:"size:32" json:"title"`
        Abstract     string `json:"abstract"`
        Content      string `json:"content"`
        Cover        string `json:"cover"`
        LookCount    int    `json:"look_count"`
        Likes        int    `json:"likes"`
        CommentCount int    `json:"comment_count"`
        CollectCount int    `json:"collect_count"`
        OpenComment  bool   `json:"open_comment"`
        // Status       int8   `json:"status"`
}

type CommitArticle struct {
        Title       string `gorm:"size:32" json:"title"`
        Abstract    string `json:"abstract"`
        Content     string `json:"content"`
        Cover       string `json:"cover"`
        OpenComment bool   `json:"open_comment"`
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
        CreatedAt    time.Time `json:"created_at"`
}
