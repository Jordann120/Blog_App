package routes

import (
	"BLOG_APP/controllers"
	"BLOG_APP/middleware"

	"github.com/gin-gonic/gin"
)

func ArticleRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		articles := api.Group("/articles")
		{
			articles.GET("", controllers.ListArticles)
			articles.GET("/:id", controllers.GetArticle)

			// Routes protégées
			auth := articles.Use(middleware.AuthMiddleware())
			{
				auth.POST("", controllers.CreateArticle)
				auth.PUT("/:id", controllers.UpdateArticle)
				auth.DELETE("/:id", controllers.DeleteArticle)

				// Nouvelles routes
				auth.POST("/:articleid/comment", controllers.AddComment)
				auth.POST("/:articleid/like", controllers.LikeArticle)
				auth.POST("/:articleid/dislike", controllers.DislikeArticle)
			}
		}
	}
}
