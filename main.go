package main

import (
	"countries-api/internal/database"
	"countries-api/internal/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"log"
)

func setupRoutes(app *fiber.App) {

	app.Get("/api/countries", routes.GetCountries)
	app.Get("/api/countries/paginated", routes.GetPaginatedCountries)

}

func main() {
	database.ConnectDB("./countries.db")

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Fiber API is running")
	})

	setupRoutes(app)

	log.Fatal(app.Listen(":3030"))

}
