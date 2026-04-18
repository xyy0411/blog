package router

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/controllers"
	"github.com/xyy0411/blog/controllers/matching"
	"github.com/xyy0411/blog/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(middlewares.CheckIP())
	api := r.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.GET("/setName", controllers.SetName, middlewares.AuthMiddlewares())
	}

	article := api.Group("/article")
	article.GET("/:article_id", controllers.ShowArticle)
	article.GET("/all", controllers.GetAllArticles)
	article.Use(middlewares.AuthMiddlewares())
	{
		article.POST("/comment/:article_id", controllers.PublishArticleComment)
		article.DELETE("/:article_id/comment/:comment_id", controllers.DeleteArticleComment)
		article.POST("", controllers.CreateArticle)
		article.PUT("/:article_id", controllers.UpdateArticle)
	}

	matched := api.Group("/matching")
	{
		matched.POST("/profile", matching.CreateMatchingProfile)
		matched.PUT("/profile/:user_id", matching.UpdateProfileName)
		matched.GET("/profile/:user_id", matching.GetMatchingProfile)
		matched.GET("/profile/:user_id/software", matching.GetMatchingSoftwareList)
		matched.GET("/profile/:user_id/block-user", matching.GetMatchingBlockUserList)
		matched.GET("/profile/:user_id/expire", matching.GetMatchingExpire)
		matched.PATCH("/profile/:user_id/expire", matching.UpdateMatchingExpire)
		matched.POST("/profile/:user_id/software", matching.AddMatchingSoftware)
		matched.DELETE("/profile/:user_id/software", matching.RemoveMatchingSoftware)
		matched.POST("/profile/:user_id/block-user", matching.AddMatchingBlockUser)
		matched.DELETE("/profile/:user_id/block-user/:target_user_id", matching.RemoveMatchingBlockUser)

		matched.GET("/record/all", matching.GetAllMatchingRecords)
		matched.GET("/record/today", matching.GetTodayMatchingRecords)
		matched.GET("/record/week", matching.GetThisWeekMatchingRecords)

		matched.GET("/status/:user_id", matching.LookMatchingStatus)
		matched.GET("/:user_id", matching.HandleMatching)
		matched.DELETE("/:user_id", matching.QuitMatching)
		matched.GET("/person", matching.GetMatchingPerson)
	}

	registerWebRoutes(r)

	return r
}

func registerWebRoutes(r *gin.Engine) {
	distDir := filepath.Join("web", "dist")
	indexPath := filepath.Join(distDir, "index.html")

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"message": "接口不存在，请检查请求路径是否正确。"})
			return
		}

		requestPath := strings.TrimPrefix(path.Clean(c.Request.URL.Path), "/")
		if requestPath != "" && requestPath != "." {
			assetPath := filepath.Join(distDir, filepath.FromSlash(requestPath))
			if fileExists(assetPath) {
				c.File(assetPath)
				return
			}
		}

		if fileExists(indexPath) {
			c.File(indexPath)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "页面不存在，请检查路径是否正确；如果是前端页面，请先在 web 目录执行 npm run build。",
		})
	})
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
