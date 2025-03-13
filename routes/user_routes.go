package routes

import (
	"BLOG_APP/controllers"
	"BLOG_APP/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/users", controllers.Register)
		api.POST("/users/login", controllers.Login)

		auth := api.Group("/user")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/", controllers.GetProfile)
			auth.PUT("/", controllers.UpdateUser)
		}
	}
}
