package controller

import (
	"github.com/feed-me/database"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	userSessionField string = "user_id"
)

type UserController struct {
	DB    database.DBTX
	Store *session.Store
}

type UserLogin struct {
	Username string
	Password string
}

func (ctrl *UserController) UserLoggedMiddleware(c *fiber.Ctx) error {
	sess, err := ctrl.Store.Get(c)
	if err != nil {
		return fiber.ErrBadRequest
	}
	sessuserid := sess.Get(userSessionField)
	if sessuserid == nil {
		return fiber.ErrForbidden
	}
	userid, ok := sessuserid.(int32)
	if !ok {
		return fiber.ErrInternalServerError
	}
	queries := database.New(ctrl.DB)
	user, err := queries.GetUserById(c.Context(), userid)
	if err != nil {
		return fiber.ErrForbidden
	}
	c.Locals("user", user)
	return c.Next()
}

func (ctrl *UserController) Login(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	var userlogin UserLogin

	if err := c.BodyParser(&userlogin); err != nil {
		return err
	}
	sess, err := ctrl.Store.Get(c)
	if err != nil {
		return err
	}
	queries := database.New(ctrl.DB)

	user := utils.GetUser(c.Context(), queries, userlogin.Username, userlogin.Password)

	if user == nil {
		return fiber.NewError(fiber.StatusForbidden, "usuário ou senha incorretos")
	}
	sess.Set(userSessionField, user.ID)
	if err := sess.Save(); err != nil {
		return err
	}
	ctxlog.Infof("%s logged in", user.Name)
	return c.JSON(fiber.Map{"name": user.Name})
}

func (ctrl *UserController) Logout(c *fiber.Ctx) error {
	ctxlog := log.WithContext(c.Context())
	sess, err := ctrl.Store.Get(c)
	if err != nil {
		return err
	}
	sess.Delete(userSessionField)
	user, _ := c.Locals("userid").(database.User)
	ctxlog.Infof("%s logged out", user.Name)
	return c.JSON(fiber.Map{"message": "success"})
}
