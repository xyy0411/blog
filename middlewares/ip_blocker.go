package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/resp"
	"net/http"
)

func CheckIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if err := global.DB.Where("ip = ?", ip).Error; err != nil {
			resp.Error(ctx, http.StatusForbidden, "知道别人公网就感觉自己牛逼了？")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
