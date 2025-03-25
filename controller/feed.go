package controller

import (
	"database/sql"
	"errors"
	"time"

	"github.com/feed-me/database"
	"github.com/feed-me/types"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type FeedController struct {
	DB *sql.DB
}

type CreateFeedRequest struct {
	Name string `json:"name"`
	EditFeedRequest
	FeedType string `json:"feed_type"`
}

type EditFeedRequest struct {
	Description *string `json:"description"`
	IsPublic    *bool   `json:"is_public"`
}

func (ctrl *FeedController) ListFeeds(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())

	queries := database.New(ctrl.DB)
	feeds, err := queries.ListFeeds(c.Context())
	if err == sql.ErrNoRows {
		ctxlog.Warn("there are no feeds")
		return fiber.ErrNotFound
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	return c.JSON(feeds)
}

func (ctrl *FeedController) CreateFeed(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	var req CreateFeedRequest
	if err := c.BodyParser(&req); err != nil {
		ctxlog.Warnf("create feed parse error: %v", err)
		ctxlog.Warnf("req: %v", req)
		return fiber.ErrBadRequest
	}
	feed_type := database.FeedsType(req.FeedType)
	if !feed_type.Valid() {
		ctxlog.Warnf("invalid feed type: %s", feed_type)
		ctxlog.Warnf("req: %v", req)
		return fiber.NewError(fiber.StatusBadRequest, "invalid feed type")
	}
	queries := database.New(ctrl.DB)
	result, err := queries.CreateFeed(c.Context(), req.Name, utils.ConvertToSqlNull(req.Description), utils.BoolColapse(req.IsPublic, true), feed_type)
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	id, err := result.LastInsertId()
	feed, err2 := queries.GetFeedById(c.Context(), int32(id))
	err = errors.Join(err, err2)
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	return c.JSON(feed)
}

func (ctrl *FeedController) EditFeed(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	var req EditFeedRequest
	id, err := c.ParamsInt("feed")
	if err != nil {
		ctxlog.Debugf("parameter feed error: %v", err)
		return fiber.ErrBadRequest
	}
	queries := database.New(ctrl.DB)

	feed, err := queries.GetFeedById(c.Context(), int32(id))
	if err == sql.ErrNoRows {
		return fiber.ErrNotFound
	}
	if err != nil {
		ctxlog.Debugf("database error: %v", err)
		return fiber.ErrBadRequest
	}

	if req.Description != nil {
		feed.Description.String = *req.Description
		feed.Description.Valid = true
	}

	err = queries.EditFeedById(c.Context(), feed.Description, utils.BoolColapse(req.IsPublic, feed.IsPublic), feed.ID)
	if err != nil {
		ctxlog.Debugf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	return c.JSON(types.JsonMessageId{Message: "feed modified", ID: int64(id)})
}

func (ctrl *FeedController) DeleteFeed(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())

	id, err := c.ParamsInt("feed")
	if err != nil {
		ctxlog.Debugf("parameter feed error: %v", err)
		return fiber.ErrBadRequest
	}
	queries := database.New(ctrl.DB)
	err = queries.RemoveFeedById(c.Context(), int32(id))
	if err == sql.ErrNoRows {
		return fiber.ErrNotFound
	}
	if err != nil {
		ctxlog.Debugf("database error: %v", err)
		return fiber.ErrBadRequest
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
		entries, err := queries.GetIPEnabledEntries(c.Context(), feed.ID, now)
		if err != nil {
			return err
		}
		for _, ip := range entries {
			c.Writef("%s\n", ip.String())
		}
	} else {
		var entries []string
		if feed.Type == database.FeedsTypeDomain {
			entries, err = queries.GetDomainEnabledEntries(c.Context(), feed.ID, now)
		} else {
			entries, err = queries.GetURLEnabledEntries(c.Context(), feed.ID, now)
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
