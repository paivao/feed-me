package controller

import (
	"database/sql"
	"errors"
	"time"

	"github.com/feed-me/database"
	"github.com/feed-me/types"
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
	result, err := queries.CreateFeed(c.Context(), req)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	return c.JSON(types.JsonMessageId{Message: "feed created successfully", ID: id})
}

func (ctrl *FeedController) EditFeed(c *fiber.Ctx) error {
	var req database.EditFeedByIdParams
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	req.ID = int32(id)
	queries := database.New(ctrl.DB)
	result, err := queries.EditFeedById(c.Context(), req)
	rows, err2 := result.RowsAffected()
	err = errors.Join(err, err2)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fiber.ErrNotFound
	}
	return c.JSON(types.JsonMessageId{Message: "Feed modified", ID: int64(id)})
}

func (ctrl *FeedController) DeleteFeed(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	queries := database.New(ctrl.DB)
	result, err := queries.RemoveFeedById(c.Context(), int32(id))
	rows, err2 := result.RowsAffected()
	err = errors.Join(err, err2)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fiber.ErrNotFound
	}
	return c.JSON(types.JsonMessageId{Message: "Feed removed", ID: int64(id)})
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
