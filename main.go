package main

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/feed-me/controller"
	"github.com/feed-me/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	log_levels = map[string]log.Level{
		"debug": log.LevelDebug,
		"error": log.LevelError,
		"fatal": log.LevelFatal,
		"info":  log.LevelInfo,
		"panic": log.LevelPanic,
		"trace": log.LevelTrace,
		"warn":  log.LevelWarn,
	}
)

func main() {
	//go:embed sql/schema/*
	var dbMigrations embed.FS

	conf, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatalf("could not load configuration: %v\n", err)
	}

	db, err := conf.ConnectDB()
	if err != nil {
		log.Fatalf("could not connect to database: %v\n", err)
	}

	log.Info("Managing migrations")

	accessLog, err := os.OpenFile(conf.Log.Access, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Fatalf("error opening access log file: %v", err)
	}

	level, ok := log_levels[conf.Log.Level]
	if !ok {
		log.Fatalf("log level not defined: %s", conf.Log.Level)
	}

	systemLog, err := os.OpenFile(conf.Log.System, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Fatalf("error opening access log file: %v", err)
	}

	log.SetLevel(level)
	log.SetOutput(systemLog)

	// Fiber instance
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${method} \"${url}\" ${protocol} ${status} ${bytesSent} \"${referer}\" \"${ua}\" ${error}\n",
		TimeFormat: time.RFC3339,
		Output:     accessLog,
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
	feedGroup.Get("", feedController.ListFeeds)
	feedGroup.Put("", feedController.CreateFeed)
	feedGroup.Post("/:id", feedController.EditFeed)
	feedGroup.Delete("/:id", feedController.DeleteFeed)

	entryGroup := api.Group("/entry", userController.UserLoggedMiddleware)
	entryGroup.Get("/:type/:feed_id/", entryController.ListEntries)
	entryGroup.Put("/:type/:feed_id/", entryController.AddEntry)
	entryGroup.Post("/:type/:feed_id/:id", entryController.EditEntry)
	entryGroup.Delete("/:type/:feed_id/:id", entryController.DeleteEntry)

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
	return ctx.Status(code).JSON(types.NewJsonError(err))
}
