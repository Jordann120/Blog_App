package middleware

import (
	"BLOG_APP/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header requis"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Ajoute les informations de l'utilisateur au contexte
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Avant la requête
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Après la requête
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			// Log les erreurs
			c.Error(nil).SetMeta(gin.H{
				"path":   path,
				"method": method,
				"status": statusCode,
			})
		}
	}
}

// ErrorHandler middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Gestion des erreurs après l'exécution de la requête
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Log l'erreur
				if meta, ok := e.Meta.(gin.H); ok {
					// Vous pouvez personnaliser la gestion des erreurs ici
					c.JSON(meta["status"].(int), gin.H{
						"error": e.Error(),
						"meta":  meta,
					})
				}
			}
		}
	}
}
