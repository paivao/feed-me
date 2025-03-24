package controller

import (
	"database/sql"
	"net/url"

	"github.com/feed-me/database"
	"github.com/feed-me/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type EntryController struct {
	DB *sql.DB
}

func (ctrl *EntryController) ListEntries(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())

	// Check feed type
	feedType := database.FeedsType(c.Params("type"))
	if !feedType.Valid() {
		ctxlog.Warnf("incorrect feed type: %v", feedType)
		return fiber.NewError(fiber.StatusBadRequest, "incorrect feed type")
	}

	// Check if limit and offset are included
	size := c.QueryInt("quantity")
	offset := c.QueryInt("offset")
	if size < 0 || offset < 0 {
		ctxlog.Warnf("invalid size (%d) and offset (%s)", size, offset)
		return fiber.ErrBadRequest
	}
	if size > 0 {
		offset *= size
	}

	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByNameAndType(c.Context(), c.Params("name"), feedType)
	if err == sql.ErrNoRows {
		ctxlog.Warnf("feed %s of type %s not found", c.Params("name"), feedType)
		return fiber.NewError(fiber.StatusNotFound, "feed not found")
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}

	var entries interface{}
	switch feedType {
	case database.FeedsTypeIp:
		if size > 0 {
			entries, err = queries.ListIPEntriesWindow(c.Context(), feed.ID, int32(size), int32(offset))
		} else {
			entries, err = queries.ListIPEntries(c.Context(), feed.ID)
		}
	case database.FeedsTypeDomain:
		if size > 0 {
			entries, err = queries.ListDomainEntriesWindow(c.Context(), feed.ID, int32(size), int32(offset))
		} else {
			entries, err = queries.ListDomainEntries(c.Context(), feed.ID)
		}
	case database.FeedsTypeUrl:
		if size > 0 {
			entries, err = queries.ListURLEntriesWindow(c.Context(), feed.ID, int32(size), int32(offset))
		} else {
			entries, err = queries.ListURLEntries(c.Context(), feed.ID)
		}
	}
	if err == sql.ErrNoRows {
		ctxlog.Warnf("empty feed")
		return fiber.NewError(fiber.StatusNotFound, "empty feed or outside limits")
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "could not retrive entries")
	}
	ctxlog.Debugf("%s retrieved by %s [%s]", feedType, c.Locals("userid"), c.Context().RemoteIP().String())
	return c.JSON(entries)
}

func (ctrl *EntryController) AddEntry(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	feedType := database.FeedsType(c.Params("type"))
	if !feedType.Valid() {
		ctxlog.Warnf("incorrect feed type: %v", feedType)
		return fiber.NewError(fiber.StatusBadRequest, "incorrect feed type")
	}
	var req struct {
		value       string
		comment     sql.NullString
		valid_until sql.NullTime
	}
	if err := c.BodyParser(&req); err != nil {
		ctxlog.Warnf("incorrect body: %v", err)
		return fiber.ErrBadRequest
	}
	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByNameAndType(c.Context(), c.Params("name"), feedType)
	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusNotFound, "feed not found")
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	var result sql.Result
	switch feedType {
	case database.FeedsTypeIp:
		var my_net types.MyNet
		err = my_net.UnmarshalJSON([]byte(req.value))
		if err != nil {
			ctxlog.Debugf("could not convert to ip: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "incorrect ip/network")
		}
		result, err = queries.InsertIPEntry(c.Context(), my_net, req.comment, req.valid_until, feed.ID)
	case database.FeedsTypeDomain:

		result, err = queries.InsertDomainEntry(c.Context(), req.value, req.comment, req.valid_until, feed.ID)
	case database.FeedsTypeUrl:
		_, err := url.Parse(req.value)
		if err != nil {
			ctxlog.Warnf("error parsing url: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "could not parse url")
		}
		result, err = queries.InsertURLEntry(c.Context(), req.value, req.comment, req.valid_until, feed.ID)
	}
	if err != nil {
		ctxlog.Warnf("error in creating new entry: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "could not create new entry")
	}
	id, err := result.LastInsertId()
	ctxlog.Infof("%s added by %s [%s]: %s", feedType, c.Locals("userid"), c.Context().RemoteIP().String(), req.value)
	return c.JSON(fiber.Map{"message": "entry added", "id": id})
}
