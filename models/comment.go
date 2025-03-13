package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Body      string  `gorm:"not null"`
	UserID    uint    `gorm:"not null"`
	ArticleID uint    `gorm:"not null"`
	Author    User    `gorm:"foreignKey:UserID"`
	Article   Article `gorm:"foreignKey:ArticleID"`
}
