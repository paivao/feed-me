package controller

import (
	"database/sql"

	"github.com/feed-me/database"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	userSessionField string = "user_id"
)

var defaultHash string

type UserController struct {
	DB    *sql.DB
	Store *session.Store
}

type UserLogin struct {
	Username string
	Password string
}

func init() {
	hash, err := utils.PasswordHash("")
	if err != nil {
		panic(err)
	}
	defaultHash = hash
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
	var userlogin UserLogin
	if err := c.BodyParser(&userlogin); err != nil {
		return err
	}
	sess, err := ctrl.Store.Get(c)
	if err != nil {
		return err
	}
	queries := database.New(ctrl.DB)
	user, err := queries.GetUserByName(c.Context(), userlogin.Username)

	if err == sql.ErrNoRows {
		user.PasswordHash = defaultHash
	} else if err != nil {
		return err
	}
	if !utils.PasswordVerify(userlogin.Password, user.PasswordHash) {
		return fiber.NewError(fiber.StatusForbidden, "usuário ou senha incorretos")
	}
	sess.Set(userSessionField, user.ID)
	if err := sess.Save(); err != nil {
		return err
	}
	return c.JSON(utils.NewMessage("logado com sucesso"))
}
