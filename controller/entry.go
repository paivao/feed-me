package controller

import (
	"database/sql"

	"github.com/feed-me/database"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
)

type EntryController struct {
	DB *sql.DB
}

func (ctrl EntryController) ListIPEntries(c *fiber.Ctx) error {
	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByName(c.Context(), c.Params("name"))
	if err != nil {
		return err
	}
	if feed.Type != database.FeedsTypeIp {
		return fiber.NewError(fiber.StatusBadRequest, "tipo incorreto de feed")
	}
	entries, err := queries.ListIPEntries(c.Context(), feed.ID)
	if err != nil {
		return err
	}
	return c.JSON(entries)
}

func (ctrl EntryController) AddIPEntry(c *fiber.Ctx) error {
	var req database.InsertIPEntryParams
	if err := c.BodyParser(req); err != nil {
		return err
	}
	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByName(c.Context(), c.Params("name"))
	if err != nil {
		return err
	}
	if feed.Type != database.FeedsTypeIp {
		return fiber.NewError(fiber.StatusBadRequest, "tipo incorreto de feed")
	}
	req.FeedID = feed.ID
	err = queries.InsertIPEntry(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(utils.NewMessage("ip added"))
}
