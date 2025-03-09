package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xyy0411/blog/server/controllers"
	"github.com/xyy0411/blog/server/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.GET("/setName", controllers.SetName, middlewares.AuthMiddlewares())
	}
	api := r.Group("/api")
	api.GET("/article/:article_id", controllers.ShowArticle)
	// 在get请求后添加一个中间件,用来检查token
	api.Use(middlewares.AuthMiddlewares())
	{
		api.POST("/article/comment/:article_id", controllers.PublishArticleComment)
		api.DELETE("/articles/:article_id/comments/:comment_id", controllers.DeleteArticleComment)
		api.POST("/article", controllers.CreateArticle)
	}
	return r
}
