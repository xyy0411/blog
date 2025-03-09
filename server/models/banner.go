package models

// Banner 首页的图片
type Banner struct {
	Model
	Cover string `json:"cover"` // 图片链接
	Href  string `json:"href"`  // 跳转链接
}
