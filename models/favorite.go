package models

type Favorite struct {
	UserID    uint `gorm:"primaryKey"`
	ArticleID uint `gorm:"primaryKey"`
}
