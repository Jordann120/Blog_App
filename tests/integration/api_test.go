package integration

import (
	"BLOG_APP/database"
	"BLOG_APP/routes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	// Configuration de la base de données de test
	database.InitDB()
	database.DB.Exec("TRUNCATE users CASCADE")
}

func TestUserFlow(t *testing.T) {
	setupTestDB()

	r := gin.Default()
	routes.UserRoutes(r)
	routes.ArticleRoutes(r)

	// Test d'inscription
	t.Run("register and login flow", func(t *testing.T) {
		// Inscription
		registerBody := map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(registerBody)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var registerResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &registerResponse)
		token := registerResponse["token"].(string)

		// Test de connexion
		loginBody := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ = json.Marshal(loginBody)

		req, _ = http.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Test de création d'article
		articleBody := map[string]interface{}{
			"title": "Test Article",
			"body":  "This is a test article",
		}
		body, _ = json.Marshal(articleBody)

		req, _ = http.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
