package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"BLOG_APP/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Configuration du test
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/users", controllers.Register)

	// Cas de test : inscription réussie
	t.Run("successful registration", func(t *testing.T) {
		reqBody := controllers.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "user")
		assert.Contains(t, response, "token")
	})

	// Cas de test : email invalide
	t.Run("invalid email", func(t *testing.T) {
		reqBody := controllers.RegisterRequest{
			Username: "testuser",
			Email:    "invalid-email",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
