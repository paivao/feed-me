package controller

import (
	"net"
	"time"

	"github.com/feed-me/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PrintFeeds(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var feed model.IPFeed
		err := db.Preload("Entries").Where("name = ?", c.Params("name")).Take(&feed).Error
		if err != nil {
			return err
		}
		now := time.Now()
		for _, ip := range feed.Entries {
			if !ip.Enabled {
				continue
			}
			if ip.ValidUntil != nil && ip.ValidUntil.Before(now) {
				continue
			}
			c.Writef("%s\n", (*net.IPNet)(&ip.Network).String())
		}
		return nil
	}
}
