package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maan78612/car_rental_dynamic/pkg/controllers"
)

func CallRoutes(app *fiber.App) {
	AuthRoutes(app)
	UserRoute(app)
	BookingRoutes(app)
}

func AuthRoutes(app *fiber.App) {
	app.Post("/user/signUp", controllers.SignUp)
	app.Post("/user/login", controllers.Login)

}

func UserRoute(app *fiber.App) {
	// app.Use(middleWare.Authenticate)
	app.Get("/user/:user_id", controllers.GetAUser)
	app.Put("/user/:user_id", controllers.EditAUser)
	app.Delete("/user/:user_id", controllers.DeleteAUser)
	app.Get("/user", controllers.GetAllUsers)

}

func BookingRoutes(app *fiber.App) {
	app.Post("/booking", controllers.Createbooking)
	app.Get("/booking/:booking_id", controllers.GetBooking)
}
