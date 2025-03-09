package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/server/global"
	"github.com/xyy0411/blog/server/models"
	"github.com/xyy0411/blog/server/utils"
	"net/http"
)

func AuthMiddlewares() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权",
			})
			ctx.Abort()
			return
		}
		uid, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token错误",
			})
			ctx.Abort()
			return
		}
		var user models.User
		if err := global.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询用户出错",
			})
		}

		ctx.Set("uid", uid)
		ctx.Next()
	}
}
