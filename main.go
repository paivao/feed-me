package main

import (
	"encoding/base64"
	"errors"
	"fmt"
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

	feedGroup := api.Group("/feed", userController.UserLoggedMiddleware)
	feedGroup.Put("/create", feedController.CreateFeed)
	feedGroup.Get("/list", feedController.ListFeeds)

	entryGroup := api.Group("/entry", userController.UserLoggedMiddleware)
	entryGroup.Get("/ip/:name/", entryController.ListIPEntries)
	entryGroup.Get("/ip/:name/create", entryController.AddIPEntry)

	app.Mount("/api", api)

	// Static file server
	app.Static("/", "./static")

	// Start server
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", conf.Host, conf.Port)))
}

func exportBasicAuth(c *fiber.Ctx) error {
	auth := c.Get(fiber.HeaderAuthorization)
	if !strings.HasPrefix(auth, "basic ") {
		return c.Next()
	}
	raw, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return fiber.ErrBadRequest
	}
	userpass := string(raw)
	index := strings.Index(userpass, ":")
	if index == -1 {
		return fiber.ErrBadRequest
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
