package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/controllers"
	"github.com/xyy0411/blog/controllers/matching"
	"github.com/xyy0411/blog/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 全局屏蔽IP
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
	// 在get请求后添加一个中间件,用来检查token
	article.Use(middlewares.AuthMiddlewares())
	{
		article.POST("/comment/:article_id", controllers.PublishArticleComment)
		article.DELETE("/:article_id/comment/:comment_id", controllers.DeleteArticleComment)
		article.POST("", controllers.CreateArticle)
		article.PUT("/:article_id", controllers.UpdateArticle)
	}
	matched := api.Group("/matching")
	{
		matched.GET("status/:user_id", matching.LookMatchingStatus)
		matched.GET("/:user_id", matching.HandleMatching)
		matched.DELETE("/:user_id", matching.QuitMatching)
		matched.GET("person", matching.GetMatchingPerson)
	}
	return r
}
