package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/maan78612/car_rental_dynamic/pkg/configs"
	"github.com/maan78612/car_rental_dynamic/pkg/routes"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.CallRoutes(app) //add this

	port := configs.EnvGetPort()

	log.Fatal(app.Listen(port))

}
