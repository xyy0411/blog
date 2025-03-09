package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/server/global"
	"github.com/xyy0411/blog/server/models"
	"github.com/xyy0411/blog/server/resp"
	"github.com/xyy0411/blog/server/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func PublishArticleComment(ctx *gin.Context) {
	var input struct {
		Content string `json:"content"`
	}

	id := ctx.Param("article_id")
	if id == "" {
		resp.Error(ctx, http.StatusBadRequest, "缺少 article_id")
		return
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	uid := ctx.MustGet("uid")

	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err)
		return
	}

	var user models.User
	if err := global.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusNotFound, "该用户可能不存在")
		return
	}

	comment := models.Comment{
		UserID:    user.ID,
		Content:   input.Content,
		ArticleID: uint(articleID),
		IP:        ctx.ClientIP(),
	}

	if err := global.DB.Create(&comment).Error; err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusInternalServerError, "发送评论失败")
		return
	}

	var a models.Article
	result := global.DB.Model(&a).Where("id = ?", articleID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	if result.Error != nil {
		global.Logger.Error(result.Error)
		resp.Error(ctx, http.StatusInternalServerError, "文章评论累加失败")
		return
	}
	if result.RowsAffected == 0 {
		global.Logger.Warn("没有找到对应的文章，更新失败")
		resp.Error(ctx, http.StatusNotFound, "文章不存在")
		return
	}

	resp.OK(ctx, "发送评论成功!", nil)
}

func DeleteArticleComment(ctx *gin.Context) {

	artID := ctx.Param("article_id")
	if artID == "" {
		resp.Error(ctx, http.StatusBadRequest, "缺少 article_id")
		return
	}

	comID := ctx.Param("comment_id")
	if comID == "" {
		resp.Error(ctx, http.StatusBadRequest, "缺少 comment_id")
		return
	}

	uid := ctx.MustGet("uid")

	if err := global.DB.Where("article_id = ? AND id = ? AND user_id = ?", artID, comID, uid).Delete(&models.Comment{}).Error; err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusInternalServerError, "删除评论失败")
		return
	}

	if err := global.DB.Model(&models.Article{}).Where("id = ?", artID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)); err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusInternalServerError, "文章评论递减失败")
		return
	}

	resp.OK(ctx, "删除评论成功!", nil)
}

func ShowArticle(ctx *gin.Context) {

	articleID := ctx.Param("article_id")
	if articleID == "" {
		resp.Error(ctx, http.StatusBadRequest, "缺少 article_id")
		return
	}

	var article models.Article
	if err := global.DB.Where("id = ?", articleID).First(&article).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	var comments []*models.Comment
	global.DB.Where("article_id = ?", articleID).Find(&comments)

	data := utils.StructToMap(models.CommentWithArticle{
		Article:  article,
		Comments: comments,
	})

	resp.OK(ctx, article.Title, data, "title")
}

func CreateArticle(ctx *gin.Context) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "错误的请求",
		})
		return
	}

	uid := ctx.MustGet("uid")

	var user models.User
	if err := global.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		global.Logger.Error(err)
		resp.Error(ctx, http.StatusInternalServerError, "用户可能不存在")
		return
	}

	article := models.Article{
		UserID:  user.ID,
		Title:   input.Title,
		Content: input.Content,
	}

	if err := global.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建文章错误",
		})
		global.Logger.Error("创建文章错误:", err)
		return
	}

	resp.OK(ctx, "创建文章成功", nil)
}
