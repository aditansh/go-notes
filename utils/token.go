package utils

import (
	"time"

	config "github.com/aditansh/go-notes/config"
	"github.com/aditansh/go-notes/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func generateToken(userID uuid.UUID, secret string, expiry time.Duration, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID.String(),
		"exp":       time.Now().Add(expiry).Unix(),
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt.Unix(),
	})

	return token.SignedString([]byte(secret))

}

func GenerateAccessToken(user *models.User) (string, error) {
	config, _ := config.LoadEnvVariables(".")
	return generateToken(user.ID, config.AccessTokenSecret, config.AccessTokenExpiry, user)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	config, _ := config.LoadEnvVariables(".")
	return generateToken(user.ID, config.RefreshTokenSecret, config.RefreshTokenExpiry, user)
}
