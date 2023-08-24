package services

import (
	"strings"
	"time"

	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(note *models.CreateNoteSchema, userID uuid.UUID) error {
	now := time.Now()

	noteModel := models.Note{
		Title:     note.Title,
		Content:   note.Content,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := database.DB.Create(&noteModel)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return fiber.NewError(fiber.StatusConflict, "Title already exists, please use another title")
	} else if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func UpdateNote(payload *models.UpdateNoteSchema, note *models.Note) error {

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}
	if payload.Category != "" {
		updates["category"] = payload.Category
	}
	if payload.Published != nil {
		updates["published"] = payload.Published
	}
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&note).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func FindNoteByID(noteId string) (models.Note, error) {
	var note models.Note
	result := database.DB.Where("id = ?", noteId).First(&note)
	if err := result.Error; err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func FindAllNotes() ([]models.Note, error) {
	var notes []models.Note
	result := database.DB.Find(&notes)
	if result.Error != nil {
		return []models.Note{}, result.Error
	}

	return notes, nil
}
