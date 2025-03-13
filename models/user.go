package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Bio       string
	Image     string
	Articles  []Article
	Comments  []Comment
	Favorites []Article `gorm:"many2many:user_favorites;"`
	Followers []User    `gorm:"many2many:follows;joinForeignKey:FollowingID;joinReferences:FollowerID"`
	Following []User    `gorm:"many2many:follows;joinForeignKey:FollowerID;joinReferences:FollowingID"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
