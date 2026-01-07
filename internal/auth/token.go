package token

import (
	"enactus/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSecret struct {
	Secret []byte
}

func (jwtSecret *JWTSecret) GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to convert jwt to string: %w", err)
	}

	return tokenStr, nil
}

func (jwtSecret *JWTSecret) ValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret.Secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
