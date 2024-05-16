package router

import (
	"eniqilo-store/handler"
	"eniqilo-store/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/v1")

	// User
	user := api.Group("/user")

	// IT
	it := user.Group("/it")
	it.Post("/register", handler.RegisterUserIT)
	it.Post("/login", handler.LoginUserIT)

	// Nurse
	nurse := user.Group("/nurse")
	nurse.Get("/", middleware.Auth(), handler.GetUsers)
	nurse.Post("/login", handler.LoginUserNurse)
	nurse.Post("/register", middleware.Auth(), handler.RegisterUserNurse)
	nurse.Post("/:id/access", middleware.Auth(), handler.ChangePasswordUserNurse)
	nurse.Put("/:id", middleware.Auth(), handler.UpdateUserNurse)
	nurse.Delete("/:id", middleware.Auth(), handler.DeleteUserNurse)

	// Medical
	medical := api.Group("/medical")

	// Patient
	patient := medical.Group("/patient")
	patient.Get("/", middleware.Auth(), handler.GetPatients)
	patient.Post("/", middleware.Auth(), handler.CreatePatient)

	// Record
	record := medical.Group("/record")
	record.Get("/", middleware.Auth(), handler.GetRecords)
	record.Post("/", middleware.Auth(), handler.CreateRecord)
}
