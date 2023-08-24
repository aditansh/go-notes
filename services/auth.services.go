package services

import (
	"fmt"
	"time"

	config "github.com/aditansh/go-notes/config"
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/aditansh/go-notes/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(payload *models.RegisterUserSchema) error {

	var user models.User
	check1 := database.DB.Where("username = ?", payload.Username).First(&user)
	if check1 != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Username already exists")
	}
	check2 := database.DB.Where("email = ?", payload.Email).First(&user)
	if check2 != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	otp, _ := utils.GenerateOTP(6)
	newUser := models.User{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Verified:  false,
		Otp:       otp,
	}

	body := fmt.Sprintf("Your OTP is %s", otp)

	email, err := utils.SendEmail(payload.Email, "OTP Verification", body)
	if err != nil {
		return err
	}

	fmt.Println(email)

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func LoginUser(payload *models.LoginUserSchema) (models.AuthResponse, error) {

	var user models.User
	result := database.DB.Where("username = ?", payload.Username).First(&user)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return models.AuthResponse{}, err
	}

	if !user.Verified {
		return models.AuthResponse{}, nil
	}

	authResponse, err := GenerateAuthTokens(&user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return authResponse, nil
}

func VerifyOTP(payload *models.VerifyOTPRequest, user *models.User) error {

	if user.Verified {
		return fmt.Errorf("user already verified")
	}

	if payload.OTP != user.Otp {
		return fmt.Errorf("invalid OTP")
	}

	user.Verified = true
	user.Otp = ""
	result := database.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ResendOTP(user *models.User) error {
	otp, _ := utils.GenerateOTP(6)

	user.Otp = otp
	result := database.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	body := fmt.Sprintf("Your OTP is %s", otp)

	email, err := utils.SendEmail(user.Email, "OTP Verification", body)
	if err != nil {
		return err
	}

	fmt.Println(email)

	return nil
}

func GenerateAuthTokens(user *models.User) (models.AuthResponse, error) {
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshTokenEntry := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}

	result := database.DB.Create(&refreshTokenEntry)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Verified:     true,
	}

	return authResponse, nil
}

func RefreshAccessToken(refreshToken string) (models.AuthResponse, *fiber.Error) {
	config, err := config.LoadEnvVariables(".")
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	user, err := GetUserByRefreshToken(refreshToken)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshTokenSecret), nil
	})
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID != user.ID {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	accessToken, err := utils.GenerateAccessToken(&user)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func UpdateUser(payload *models.UpdateUserSchema, user *models.User) error {

	updates := make(map[string]interface{})
	if payload.Username != "" {
		updates["username"] = payload.Username
	}
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	return nil
}