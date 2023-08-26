package services

import (
	"github.com/aditansh/go-notes/cache"
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/gofiber/fiber/v2"
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

	id, err := cache.GetValue(refreshToken)
	// fmt.Println(id)
	// fmt.Println(err)
	if err != nil {
		return models.User{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
