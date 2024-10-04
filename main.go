package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/feed-me/controller"
	"github.com/feed-me/model"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

func networkFromCIDR(s string) model.Net {
	_, net, err := net.ParseCIDR(s)
	if err != nil {
		log.Fatal(err)
	}
	return model.Net(*net)
}

func main() {
	conf, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := conf.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting database migration...")
	err = db.AutoMigrate(&model.IPFeed{}, &model.IPEntry{}, &model.User{}, &model.Group{}, &model.Permission{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database migration done.")

	generateTestData(db)

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

	// Expose feed list
	app.Get("/feed/ip/:name", exportBasicAuth, controller.PrintFeeds(db))

	//Api
	api := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})
	api.Post("/login", controller.Login(db, store))
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

func generateTestData(db *gorm.DB) {
	test := model.IPFeed{
		Feed: model.Feed{
			Name: "teste",
		},
		Entries: []model.IPEntry{
			{
				Entry: model.Entry{
					Enabled: true,
				},
				Network: networkFromCIDR("127.0.0.1/32"),
			},
		},
	}
	hash, err := utils.PasswordHash("admin")
	if err != nil {
		log.Fatal(err)
	}
	testUser := model.User{
		Name:         "admin",
		PasswordHash: hash,
		Email:        "admin@feed.me",
	}
	db.Create(&test)
	db.Create(&testUser)
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
