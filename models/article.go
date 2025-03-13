package models

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Body        string    `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	Author      User      `gorm:"foreignKey:UserID"`
	Comments    []Comment
	FavoredBy   []User    `gorm:"many2many:user_favorites;"`
	Likes       int       `gorm:"default:0"`
}