package controller

import (
	"net"
	"time"

	"github.com/feed-me/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FeedController baseController

func NewFeedController(db *gorm.DB) FeedController {
	return FeedController{
		db: db,
	}
}

func (c FeedController) PrintFeeds(f *fiber.Ctx) error {
	var feed model.IPFeed
	err := c.db.Preload("Entries").Where("name = ?", f.Params("name")).Take(&feed).Error
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
		f.Writef("%s\n", (*net.IPNet)(&ip.Network).String())
	}
	return nil
}
