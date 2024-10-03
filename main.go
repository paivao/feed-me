package main

import (
	"encoding/base64"
	"log"
	"net"
	"strings"
	"time"
	"crypto/"

	"github.com/feed-me/controller"
	"github.com/feed-me/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	testUser := model.User{
		Name: "admin",
		PasswordHash: ,
	}
	db.Create(&test)

	feedCtrl := controller.NewFeedController(db)
	userCtrl := controller.NewUserController(db)

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
	app.Get("/feed/ip/:name", exportBasicAuth, feedCtrl.PrintFeeds)

	//Api
	api := app.Group("/api")
	api.Post("/login", userCtrl.Login)

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
	return nil
}
