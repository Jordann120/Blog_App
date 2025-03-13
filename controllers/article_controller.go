package controllers

import (
	"BLOG_APP/database"
	"BLOG_APP/models"
	"BLOG_APP/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Body        string `json:"body" binding:"required"`
}

type CommentRequest struct {
	Body string `json:"body" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	article := models.Article{
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		UserID:      userID.(uint),
	}

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": article})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if err := database.DB.Preload("Author").Preload("Comments").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Non autorisé"})
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"body":        req.Body,
	}

	if err := database.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Non autorisé"})
		return
	}

	if err := database.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article supprimé"})
}

func ListArticles(c *gin.Context) {
	var articles []models.Article
	query := database.DB.Preload("Author")

	// Filtrage par auteur
	if author := c.Query("author"); author != "" {
		query = query.Joins("JOIN users ON users.id = articles.user_id").
			Where("users.username = ?", author)
	}

	if err := query.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func AddComment(c *gin.Context) {
	articleID := c.Param("articleid")
	userID, _ := c.Get("userID")

	// Vérifier d'abord que l'article existe
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	comment := models.Comment{
		Body:      req.Body,
		ArticleID: article.ID, // Utilisation de l'ID correct
		UserID:    userID.(uint),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du commentaire"})
		return
	}

	// Charger l'auteur du commentaire pour la réponse
	database.DB.Preload("Author").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

// LikeArticle et DislikeArticle - Ajout de vérifications supplémentaires
func LikeArticle(c *gin.Context) {
	articleID := c.Param("articleid")
	userID, _ := c.Get("userID")

	// Vérifier que l'article existe
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	// Vérifier si l'utilisateur n'a pas déjà liké l'article
	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, article.ID).First(&favorite)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article déjà liké"})
		return
	}

	// Créer le like
	favorite = models.Favorite{
		UserID:    userID.(uint),
		ArticleID: article.ID,
	}

	// Utiliser une transaction pour assurer la cohérence
	tx := database.DB.Begin()
	if err := tx.Create(&favorite).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du like de l'article"})
		return
	}

	if err := tx.Model(&article).Update("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du like de l'article"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Article liké avec succès"})
}

func DislikeArticle(c *gin.Context) {
	articleID := c.Param("articleid")
	userID, _ := c.Get("userID")

	// Vérifier que l'article existe
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
		return
	}

	// Vérifier si l'utilisateur a liké l'article
	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, article.ID).First(&favorite)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article non liké"})
		return
	}

	// Utiliser une transaction pour assurer la cohérence
	tx := database.DB.Begin()
	if err := tx.Delete(&favorite).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du dislike de l'article"})
		return
	}

	if err := tx.Model(&article).Update("likes", gorm.Expr("GREATEST(likes - ?, 0)", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du dislike de l'article"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Article disliké avec succès"})
}
