package controllers

import (
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/aditansh/go-notes/services"
	"github.com/aditansh/go-notes/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateNote(c *fiber.Ctx) error {
	var payload models.CreateNoteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	userID := c.Locals("userID").(uuid.UUID)

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors})
	}

	err := services.CreateNote(&payload, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Note Created Successfully"})
}

func FindAllNotes(c *fiber.Ctx) error {
	var notes []models.Note

	notes, err := services.FindAllNotes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Notes Found Successfully",
		"data":    notes})
}

func FindNote(c *fiber.Ctx) error {
	noteId := c.Params("id")
	var note models.Note

	note, err := services.FindNoteByID(noteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Note Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Note Found Successfully",
		"data":    note})
}

func UpdateNote(c *fiber.Ctx) error {
	noteId := c.Params("id")

	var payload models.UpdateNoteSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Bad Request"})
	}

	userID := c.Locals("userID").(uuid.UUID)
	var note models.Note

	result := database.DB.Where("id = ? AND user_id = ?", noteId, userID).First(&note)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Note Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors})
	}

	err := services.UpdateNote(&payload, &note)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Note Updated Successfully"})
}

func DeleteNote(c *fiber.Ctx) error {
	noteId := c.Params("id")

	result := database.DB.Delete(&models.Note{}, "id= ?", noteId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Note Not Found"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Note Deleted Successfully"})
}
