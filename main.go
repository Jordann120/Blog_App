// main.go
package main

import (
	"BLOG_APP/database"
	"BLOG_APP/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	routes.UserRoutes(r)
	routes.ArticleRoutes(r)

	fmt.Println("🚀 Serveur lancé sur http://localhost:8080")
	r.Run(":8080")
}
