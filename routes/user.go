package routes

import (
	"github.com/aditansh/go-notes/controllers"
	"github.com/aditansh/go-notes/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {

	user := app.Group("/user")
	user.Post("/signup", controllers.SignupUser)
	user.Post("/login", controllers.LoginUser)
	user.Post("/verify", controllers.VerifyOTP)
	user.Post("/resend", controllers.ResendOTP)
	// user.Get("/findall", controllers.FindAllUsers)
	// user.Get("/find/:id", controllers.FindUser)
	user.Post("/refresh", controllers.RefreshToken)
	user.Get("/logout", middleware.ValidateToken, controllers.LogoutUser)
	user.Post("/update", middleware.ValidateToken, controllers.UpdateUser)
	user.Delete("/delete", middleware.ValidateToken, controllers.DeleteUser)
	user.Get("/me", middleware.ValidateToken, controllers.Me)
}
