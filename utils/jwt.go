package utils

import (
	"BLOG_APP/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userID uint, email string) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.LoadConfig().JWTSecret))
}

func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.LoadConfig().JWTSecret), nil
		},
	)
	
	if err != nil {
		return nil, err
	}
	
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, err
	}
	
	return claims, nil
} 