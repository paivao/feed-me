package controller

import (
	"database/sql"
	"time"

	"github.com/feed-me/database"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
)

type FeedController struct {
	DB *sql.DB
}

func (ctrl *FeedController) ListFeeds(c *fiber.Ctx) error {
	queries := database.New(ctrl.DB)
	feeds, err := queries.ListFeeds(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(feeds)
}

func (ctrl *FeedController) CreateFeed(c *fiber.Ctx) error {
	var req database.CreateFeedParams
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	queries := database.New(ctrl.DB)
	err := queries.CreateFeed(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(utils.NewMessage("feed criado com sucesso"))
}

func (ctrl *FeedController) PrintFeed(c *fiber.Ctx) error {
	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByName(c.Context(), c.Params("name"))
	if err == sql.ErrNoRows {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}
	now := sql.NullTime{Time: time.Now(), Valid: true}

	if feed.Type == database.FeedsTypeIp {
		entries, err := queries.GetIPEnabledEntries(c.Context(), database.GetIPEnabledEntriesParams{Enabled: true, FeedID: feed.ID, ValidUntil: now})
		if err != nil {
			return err
		}
		for _, ip := range entries {
			c.Writef("%s\n", ip.String())
		}
	} else {
		var entries []string
		if feed.Type == database.FeedsTypeDomain {
			entries, err = queries.GetDomainEnabledEntries(c.Context(),
				database.GetDomainEnabledEntriesParams{Enabled: true, FeedID: feed.ID, ValidUntil: now})
		} else {
			entries, err = queries.GetURLEnabledEntries(c.Context(), database.GetURLEnabledEntriesParams{Enabled: true, FeedID: feed.ID, ValidUntil: now})
		}
		if err != nil {
			return err
		}
		for _, value := range entries {
			c.Writef("%s\n", value)
		}
	}
	return nil
}
