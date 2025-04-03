package main

import (
	"countries-api/internal/database"
	"countries-api/internal/routes"
	"github.com/gofiber/fiber/v3"
	"log"
)

func setupRoutes(app *fiber.App) {

	app.Get("/api/countries", routes.GetCountries)

}

func main() {
	database.ConnectDB("./countries.db")

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Fiber API is running")
	})

	setupRoutes(app)

	log.Fatal(app.Listen(":3030"))

}
