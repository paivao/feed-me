package main

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/feed-me/controller"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
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

	store := session.New()

	// Fiber instance
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${method} \"${url}\" ${protocol} ${status} ${bytesSent} \"${referer}\" \"${ua}\" ${error}\n",
		TimeFormat: time.RFC3339,
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: conf.Key,
	}))

	feedController := controller.FeedController{DB: db}
	entryController := controller.EntryController{DB: db}
	userController := controller.UserController{DB: db, Store: store}

	// Expose feed list
	app.Get("/feed/:name", exportBasicAuth, feedController.PrintFeed)

	//Api
	api := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})
	api.Post("/login", userController.Login)

	api.Put("/feed/create", feedController.CreateFeed)
	api.Get("/feed/list", feedController.ListFeeds)

	api.Get("/entry/ip/:name/", entryController.ListIPEntries)
	api.Get("/entry/ip/:name/create", entryController.AddIPEntry)

	app.Mount("/api", api)

	// Static file server
	app.Static("/", "./static")

	// Start server
	log.Fatal(app.Listen(":3000"))
}

func exportBasicAuth(c *fiber.Ctx) error {
	auth := c.Get(fiber.HeaderAuthorization)
	if !strings.HasPrefix(auth, "basic ") {
		return nil
	}
	raw, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return nil
	}
	userpass := string(raw)
	index := strings.Index(userpass, ":")
	if index == -1 {
		return nil
	}
	c.Locals("username", userpass[:index])
	c.Locals("password", userpass[index+1:])
	return c.Next()
}

func jsonErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Return status code with error message
	return ctx.Status(code).JSON(utils.NewJsonError(err))
}
