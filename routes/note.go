package routes

import (
	"github.com/aditansh/go-notes/controllers"
	"github.com/aditansh/go-notes/middleware"
	"github.com/gofiber/fiber/v2"
)

func NoteRoutes(app *fiber.App) {

	note := app.Group("/note", middleware.ValidateToken)
	note.Post("/create", controllers.CreateNote)
	note.Get("/findall", controllers.FindAllNotes)
	note.Get("/find/:id", controllers.FindNote)
	note.Put("/update/:id", controllers.UpdateNote)
	note.Delete("/delete/:id", controllers.DeleteNote)
}
