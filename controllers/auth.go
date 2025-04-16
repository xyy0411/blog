package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	"github.com/xyy0411/blog/resp"
	"github.com/xyy0411/blog/utils"
	"gorm.io/gorm"
	"net/http"
)

func SetName(ctx *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "错误的请求",
		})
		return
	}

	name := ctx.MustGet("name").(string)
	var user models.User
	if err := global.DB.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "用户不存在")
		} else {
			resp.Error(ctx, http.StatusInternalServerError, "服务器数据库出现问题")
			global.Logger.Error(err)
		}
		return
	}

	user.Nickname = input.Name

	if err := global.DB.Save(&user).Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "保存用户数据失败")
		global.Logger.Error(err)
		return
	}

	resp.OK(ctx, "更新用户名成功！", nil)
}

func Login(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "错误的请求",
		})
		return
	}

	var dbUser models.User
	if err := global.DB.Where("nickname = ?", input.Name).First(&dbUser).Error; err != nil {
		global.Logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在的用户",
		})
		return
	}

	if !utils.CheckPassword(input.Password, dbUser.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "密码错误",
		})
		return
	}

	token, err := utils.GenerateJWT(dbUser.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	resp.OK(ctx, token, nil, "token")
}

func Register(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "错误的请求",
		})
		return
	}

	user := models.User{
		UID:      uuid.New().String(),
		Nickname: input.Name,
		Password: input.Password,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusInternalServerError, "注册失败")
		return
	}

	token, err := utils.GenerateJWT(user.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	resp.OK(ctx, token, map[string]any{"uid": user.UID}, "token")
}
