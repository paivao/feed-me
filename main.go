package main

import (
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/feed-me/controller"
	"github.com/feed-me/database"
	"github.com/feed-me/types"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jackc/pgx/v5"
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

// Embed a single file
//
//go:embed index.html
var index_page embed.FS

// Embed a directory
//
//go:embed static/*
var embed_static embed.FS

func main() {
	bg := context.Background()
	conf, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatalf("could not load configuration: %v\n", err)
	}

	db, err := conf.ConnectDB(bg)
	if err != nil {
		log.Fatalf("could not connect to database: %v\n", err)
	}
	defer db.Close(bg)

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

	store := session.New()
	log.SetLevel(level)
	log.SetOutput(systemLog)

	// Fiber instance
	fiber_app := fiber.New()
	fiber_app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${method} \"${url}\" ${protocol} ${status} ${bytesSent} \"${referer}\" \"${ua}\" ${error}\n",
		TimeFormat: time.RFC3339,
		Output:     accessLog,
	}))
	fiber_app.Use(encryptcookie.New(encryptcookie.Config{
		Key: conf.Key,
	}))

	var app fiber.Router
	if conf.BasePath == "" || conf.BasePath == "/" {
		app = fiber_app
	} else {
		app = fiber_app.Group(conf.BasePath)
	}

	feedController := controller.FeedController{DB: db}
	entryController := controller.EntryController{DB: db}
	userController := controller.UserController{DB: db, Store: store}

	// Expose feed list
	app.Get("/feed/:name", exportBasicAuth(db), feedController.PrintFeed)

	//Api
	api := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})
	api.Post("/login", userController.Login)
	api.Post("/logout", userController.Logout)
	api.Get("/whoami", userController.UserLoggedMiddleware, func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(database.User)
		if !ok {
			return fiber.ErrNotFound
		}
		return c.JSON(fiber.Map{"name": user.Name})
	})

	feedGroup := api.Group("/feed", userController.UserLoggedMiddleware)
	feedGroup.Get("", feedController.ListFeeds)
	feedGroup.Put("", feedController.CreateFeed)
	feedGroup.Post("/:feed", feedController.EditFeed)
	feedGroup.Delete("/:feed", feedController.DeleteFeed)

	entryGroup := api.Group("/entry", userController.UserLoggedMiddleware)
	entryGroup.Get("/:type/:feed/", entryController.ListEntries)
	entryGroup.Put("/:type/:feed/", entryController.AddEntry)
	entryGroup.Post("/:type/:feed/:entry", entryController.EditEntry)
	entryGroup.Delete("/:type/:feed/:entry", entryController.RemoveEntry)

	app.Mount("/api", api)

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(index_page),
	}))

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embed_static),
		PathPrefix: "static",
		Browse:     false,
	}))

	// Start server
	log.Fatal(fiber_app.Listen(fmt.Sprintf("%s:%d", conf.Host, conf.Port)))
}

func exportBasicAuth(db *pgx.Conn) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		auth := c.Get(fiber.HeaderAuthorization)
		c.Locals("user", nil)
		if !strings.HasPrefix(auth, "basic ") && !strings.HasPrefix(auth, "Basic ") {
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
		user := utils.GetUser(c.Context(), database.New(db), userpass[:index], userpass[index+1:])
		c.Locals("user", user)
		return c.Next()
	}
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
