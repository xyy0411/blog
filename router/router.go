package router

import (
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
		matched.GET("/profile/:user_id", matching.GetMatchingProfile)
		matched.GET("/profile/:user_id/software", matching.GetMatchingSoftwareList)
		matched.GET("/profile/:user_id/block-user", matching.GetMatchingBlockUserList)
		matched.GET("/profile/:user_id/expire", matching.GetMatchingExpire)
		matched.PATCH("/profile/:user_id/expire", matching.UpdateMatchingExpire)
		matched.POST("/profile/:user_id/software", matching.AddMatchingSoftware)
		matched.DELETE("/profile/:user_id/software/:software_name", matching.RemoveMatchingSoftware)
		matched.POST("/profile/:user_id/block-user", matching.AddMatchingBlockUser)
		matched.DELETE("/profile/:user_id/block-user/:target_user_id", matching.RemoveMatchingBlockUser)

		matched.GET("/status/:user_id", matching.LookMatchingStatus)
		matched.GET("/:user_id", matching.HandleMatching)
		matched.DELETE("/:user_id", matching.QuitMatching)
		matched.GET("/person", matching.GetMatchingPerson)
	}
	return r
}
