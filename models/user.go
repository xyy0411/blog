package models

// User 用户自己的一些信息
type User struct {
	Model
	UID            string `json:"uid"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`          // 头像
	Abstract       string `json:"abstract"`        // 简介
	RegisterSource int8   `json:"register_source"` // 注册来源
	Password       string `json:"-" structToMap:"ignore"`
	Email          string `json:"email"`
}
