package controller

import (
	"errors"

	"github.com/feed-me/model"
	"github.com/feed-me/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

var defaultHash string

func init() {
	hash, err := utils.PasswordHash("")
	if err != nil {
		panic(err)
	}
	defaultHash = hash
}

func Login(db *gorm.DB, store *session.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var userlogin struct {
			Username string
			Password string
		}
		if err := c.BodyParser(&userlogin); err != nil {
			return err
		}
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		var user model.User
		err = db.Where("name = ?", userlogin.Username).Take(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.PasswordHash = defaultHash
		} else if err != nil {
			return err
		}
		if !utils.PasswordVerify(userlogin.Password, user.PasswordHash) {
			return fiber.NewError(fiber.StatusForbidden, "usuário ou senha incorretos")
		}
		sess.Set("username", user.ID)
		return c.JSON(utils.NewMessage("logado com sucesso"))
	}
}
