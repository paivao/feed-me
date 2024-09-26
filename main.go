package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := conf.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Fiber instance
	app := fiber.New()
	// Expose feed list
	app.Get("/feed/:name", func(c *fiber.Ctx) error {

	})

	//Api
	api := app.Group("/api")

	// Static file server
	app.Static("/", "./static")

	// Start server
	log.Fatal(app.Listen(":3000"))
}
