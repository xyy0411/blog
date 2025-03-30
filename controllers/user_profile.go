package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	"github.com/xyy0411/blog/resp"
	"github.com/xyy0411/blog/utils"
	"gorm.io/gorm"
	"net/http"
)

func ShowUserProfile(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if uid == "" {
		resp.Error(ctx, http.StatusBadRequest, "缺少 uid")
		return
	}

	var user models.User
	if err := global.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "用户不存在")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询用户出错")
		global.Logger.Error(err)
		return
	}

	resp.OK(ctx, "查询用户成功!", utils.StructToMap(user))

}
