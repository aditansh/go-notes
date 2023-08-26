package controllers

import (
	"github.com/aditansh/go-notes/cache"
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/aditansh/go-notes/services"
	"github.com/aditansh/go-notes/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SignupUser(c *fiber.Ctx) error {
	var payload models.RegisterUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors})
	}

	err := services.SignupUser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Signup Successful. Please verify your email"})
}

func LoginUser(c *fiber.Ctx) error {
	var payload models.LoginUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	authResponse, err := services.LoginUser(&payload)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "User not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	if !authResponse.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Please verify your email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Login Successful",
		"data":    fiber.Map{"auth": authResponse}})

}

func VerifyOTP(c *fiber.Ctx) error {
	var payload models.VerifyOTPRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	user, err := services.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "User not found"})
	}

	// if user.Verified == true {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  false,
	// 		"message": "User already verified"})
	// }

	err = services.VerifyOTP(&payload, &user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	authResponse, err := services.GenerateAuthTokens(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "OTP Verified",
		"data":    fiber.Map{"auth": authResponse}})
}

func ResendOTP(c *fiber.Ctx) error {
	var payload models.ResendOTPRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	user, err := services.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"massage": "User not found"})
	}

	if user.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "User already verified"})
	}

	err = services.ResendOTP(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "OTP Sent"})
}

func RefreshToken(c *fiber.Ctx) error {

	type tokenRequest struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	var payload tokenRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	authResponse, err := services.RefreshAccessToken(payload.RefreshToken)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Token Refreshed",
		"data":    fiber.Map{"auth": authResponse}})
}

func LogoutUser(c *fiber.Ctx) error {

	type tokenRequest struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	var payload tokenRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	err := cache.DeleteValue(payload.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "User Logged Out Successfully"})
}

func UpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var payload models.UpdateUserSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors})
	}

	// var user models.User
	user, err := services.GetUserById(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	err = services.UpdateUser(&payload, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "User Updated Successfully"})
}

func DeleteUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	result := database.DB.Delete(&models.User{}, "id=?", userID)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "User Not Found"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "User Deleted Successfully"})
}

func Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// var user models.User
	user, err := services.GetUserById(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "User Found",
		"data":    user})
}
