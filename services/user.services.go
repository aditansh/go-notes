package services

import (
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/google/uuid"
)

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func GetUserById(id uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func GetUserByRefreshToken(refreshToken string) (models.User, error) {
	var token models.RefreshToken
	result := database.DB.Where("token = ?", refreshToken).First(&token)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	var user models.User
	result = database.DB.Where("id = ?", token.UserID).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
