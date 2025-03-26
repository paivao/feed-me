package controller

import (
	"database/sql"
	"net/url"
	"strconv"
	"time"

	"github.com/feed-me/database"
	"github.com/feed-me/types"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type EntryController struct {
	DB *sql.DB
}

type NewEntryRequest struct {
	Value string `json:"value"`
	commonEntryRequest
}

type EditEntryRequest struct {
	Enabled *bool `json:"enabled"`
	commonEntryRequest
}

type commonEntryRequest struct {
	Description *string    `json:"description"`
	ValidUntil  *time.Time `json:"valid_until"`
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
	feed_id, err := c.ParamsInt("feed")

	queries := database.New(ctrl.DB)
	feed, err := queries.GetFeedByIdAndType(c.Context(), int32(feed_id), feedType)
	if err == sql.ErrNoRows {
		ctxlog.Warnf("feed %s of type %s not found", c.Params("feed"), feedType)
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
	ctxlog.Debugf("%v", entries)
	return c.JSON(entries)
}

func (ctrl *EntryController) AddEntry(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	feedType := database.FeedsType(c.Params("type"))
	if !feedType.Valid() {
		ctxlog.Warnf("incorrect feed type: %v", feedType)
		return fiber.NewError(fiber.StatusBadRequest, "incorrect feed type")
	}
	var req NewEntryRequest
	if err := c.BodyParser(&req); err != nil {
		ctxlog.Warnf("incorrect body: %v", err)
		return fiber.ErrBadRequest
	}
	queries := database.New(ctrl.DB)
	feed_id, err := c.ParamsInt("feed")
	feed, err := queries.GetFeedByIdAndType(c.Context(), int32(feed_id), feedType)
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
		ctxlog.Warnf("%v", req)
		var my_net types.MyNet
		err = my_net.FromString(req.Value)
		if err != nil {
			ctxlog.Debugf("could not convert to ip: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "incorrect ip/network")
		}
		result, err = queries.InsertIPEntry(c.Context(), my_net, utils.ConvertToNullString(req.Description), utils.ConvertToNullTime(req.ValidUntil), feed.ID)
	case database.FeedsTypeDomain:

		result, err = queries.InsertDomainEntry(c.Context(), req.Value, utils.ConvertToNullString(req.Description), utils.ConvertToNullTime(req.ValidUntil), feed.ID)
	case database.FeedsTypeUrl:
		_, err := url.Parse(req.Value)
		if err != nil {
			ctxlog.Warnf("error parsing url: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "could not parse url")
		}
		result, err = queries.InsertURLEntry(c.Context(), req.Value, utils.ConvertToNullString(req.Description), utils.ConvertToNullTime(req.ValidUntil), feed.ID)
	}
	if err != nil {
		ctxlog.Warnf("error in creating new entry: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "could not create new entry")
	}
	id, err := result.LastInsertId()
	ctxlog.Infof("%s added by %s [%s]: %s", feedType, c.Locals("userid"), c.Context().RemoteIP().String(), req.Value)
	return c.JSON(fiber.Map{"message": "entry added", "id": id})
}

func (ctrl *EntryController) EditEntry(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())

	id, err := strconv.ParseInt(c.Params("entry", "0"), 10, 64)
	if err != nil {
		ctxlog.Warnf("incorrect entry id: %s, %v", c.Params("entry"), err)
		return fiber.NewError(fiber.StatusBadRequest, "entry is not integer")
	}

	feedType := database.FeedsType(c.Params("type"))
	if !feedType.Valid() {
		ctxlog.Warnf("incorrect feed type: %v", feedType)
		return fiber.NewError(fiber.StatusBadRequest, "incorrect feed type")
	}
	var req EditEntryRequest
	if err := c.BodyParser(&req); err != nil {
		ctxlog.Warnf("incorrect body: %v", err)
		return fiber.ErrBadRequest
	}
	queries := database.New(ctrl.DB)
	feed_id, err := c.ParamsInt("feed")

	feed, err := queries.GetFeedByIdAndType(c.Context(), int32(feed_id), feedType)
	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusNotFound, "feed not found")
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}
	switch feedType {
	case database.FeedsTypeIp:
		ipnet, err := queries.GetIPEntryById(c.Context(), id, feed.ID)
		if err == sql.ErrNoRows {
			return fiber.ErrNotFound
		}
		if err != nil {
			ctxlog.Warnf("database error: %v", err)
			return fiber.ErrBadRequest
		}
		if req.Description != nil {
			ipnet.Description = sql.NullString{Valid: true, String: *req.Description}
		}
		err = queries.EditIPEntryById(c.Context(), utils.BoolColapse(req.Enabled, ipnet.Enabled), ipnet.Description, utils.ConvertToNullTime(req.ValidUntil), ipnet.ID)
	case database.FeedsTypeDomain:
		domain, err := queries.GetDomainEntryById(c.Context(), id, feed.ID)
		if err == sql.ErrNoRows {
			return fiber.ErrNotFound
		}
		if err != nil {
			ctxlog.Warnf("database error: %v", err)
			return fiber.ErrBadRequest
		}
		if req.Description != nil {
			domain.Description = sql.NullString{Valid: true, String: *req.Description}
		}
		err = queries.EditDomainEntryById(c.Context(), utils.BoolColapse(req.Enabled, domain.Enabled), domain.Description, utils.ConvertToNullTime(req.ValidUntil), domain.ID)
	case database.FeedsTypeUrl:
		url, err := queries.GetURLEntryById(c.Context(), id, feed.ID)
		if err == sql.ErrNoRows {
			return fiber.ErrNotFound
		}
		if err != nil {
			ctxlog.Warnf("database error: %v", err)
			return fiber.ErrBadRequest
		}
		if req.Description != nil {
			url.Description = sql.NullString{Valid: true, String: *req.Description}
		}
		err = queries.EditDomainEntryById(c.Context(), utils.BoolColapse(req.Enabled, url.Enabled), url.Description, utils.ConvertToNullTime(req.ValidUntil), url.ID)
	}
	if err == sql.ErrNoRows {
		ctxlog.Warnf("error in removing entry: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "entry not found")
	}
	if err != nil {
		ctxlog.Warnf("error in editing entry: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "could not edit entry")
	}
	ctxlog.Infof("%s edited by %s [%s]: %d/%d", feedType, c.Locals("user"), c.Context().RemoteIP().String(), id, feed.ID)
	return c.JSON(fiber.Map{"message": "entry editted", "id": id})
}

func (ctrl *EntryController) RemoveEntry(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())

	id, err := strconv.ParseInt(c.Params("entry", "0"), 10, 64)
	if err != nil {
		ctxlog.Warnf("incorrect entry id: %s, %v", c.Params("entry"), err)
		return fiber.NewError(fiber.StatusBadRequest, "entry is not integer")
	}

	feedType := database.FeedsType(c.Params("type"))
	if !feedType.Valid() {
		ctxlog.Warnf("incorrect feed type: %v", feedType)
		return fiber.NewError(fiber.StatusBadRequest, "incorrect feed type")
	}

	queries := database.New(ctrl.DB)
	feed_id, err := c.ParamsInt("feed")

	feed, err := queries.GetFeedByIdAndType(c.Context(), int32(feed_id), feedType)
	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusNotFound, "feed not found")
	}
	if err != nil {
		ctxlog.Warnf("database error: %v", err)
		return fiber.ErrBadRequest
	}

	switch feedType {
	case database.FeedsTypeIp:
		err = queries.RemoveIPEntry(c.Context(), id, feed.ID)
	case database.FeedsTypeDomain:
		err = queries.RemoveDomainEntry(c.Context(), id, feed.ID)
	case database.FeedsTypeUrl:
		err = queries.RemoveURLEntry(c.Context(), id, feed.ID)
	}
	if err == sql.ErrNoRows {
		ctxlog.Warnf("error in removing entry: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "entry not found")
	}
	if err != nil {
		ctxlog.Warnf("error in removing entry: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "could not remove entry")
	}
	ctxlog.Infof("%s removed by %s [%s]: %d/%d", feedType, c.Locals("user"), c.Context().RemoteIP().String(), id, feed.ID)
	return c.JSON(fiber.Map{"message": "entry removed", "id": id})
}
