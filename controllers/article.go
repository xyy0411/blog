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
	"strconv"
)

func UpdateArticle(ctx *gin.Context) {
	var input models.CommitArticle

	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err)
		return
	}

	articleID := ctx.Param("article_id")

	uid := ctx.MustGet("uid")

	var user models.User
	if err := global.DB.Where("uid =?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "用户不存在")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询数据库错误")
		global.Logger.Error(err)
		return
	}

	var article models.Article
	// 使用事务处理更新操作
	tx := global.DB.Begin()
	if tx.Error != nil {
		resp.Error(ctx, http.StatusInternalServerError, "开始事务失败")
		global.Logger.Error(tx.Error)
		return
	}

	// 查询文章并检查用户权限
	if err := tx.Where("id = ? AND user_id = ?", articleID, user.ID).First(&article).Error; err != nil {
		tx.Rollback() // 回滚事务
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "文章不存在或你不是文章作者")
		} else {
			resp.Error(ctx, http.StatusInternalServerError, "查询文章出错")
			global.Logger.Error(err)
		}
		return
	}

	// 更新文章标题和内容
	article.Title = input.Title
	article.Content = input.Contents.String()

	// 保存更新后的文章
	if err := tx.Save(&article).Error; err != nil {
		tx.Rollback() // 回滚事务
		resp.Error(ctx, http.StatusInternalServerError, "更新文章失败")
		global.Logger.Error(err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "提交事务失败")
		global.Logger.Error(err)
		return
	}
	resp.OK(ctx, "更新文章成功!", nil)
}

func GetAllArticles(ctx *gin.Context) {
	var articles []*models.Article
	if err := global.DB.Find(&articles).Error; err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询文章失败")
		global.Logger.Error(err)
		return
	}
	resp.OK(ctx, "查询文章成功", utils.StructToMap(articles))
}

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

	rowsAffected := global.DB.Where("article_id = ? AND id = ? AND user_id = ?", artID, comID, user.ID).Delete(&models.Comment{})

	if rowsAffected.Error != nil {
		global.Logger.Errorf("删除评论失败: article_id=%v, comment_id=%v, user_id=%v, error=%v", artID, comID, uid, rowsAffected.Error)
		resp.Error(ctx, http.StatusInternalServerError, "删除评论失败")
		return
	}

	if rowsAffected.RowsAffected == 0 {
		global.Logger.Warnf("未找到需要删除的评论: article_id=%v, comment_id=%v, user_id=%v", artID, comID, uid)
		resp.Error(ctx, http.StatusNotFound, "评论不存在")
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
	var input models.CommitArticle
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "错误的请求",
		})
		return
	}

	if input.Title == "" || input.Contents.Len() == 0 {
		resp.Error(ctx, http.StatusBadRequest, "标题和内容不能为空")
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
		UserID:      user.ID,
		Title:       input.Title,
		Content:     input.Contents.String(),
		Cover:       input.Cover,
		OpenComment: input.OpenComment,
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
