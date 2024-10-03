package controller

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController baseController

func NewUserController(db *gorm.DB) UserController {
	return UserController{
		db: db,
	}
}

func (c UserController) Login(f *fiber.Ctx) error {
	var userlogin struct {
		Username string
		Password string
	}
	if err := f.BodyParser(&userlogin); err != nil {
		return err
	}
	f.Writef("Welcome, %s (%s)\n", userlogin.Username, userlogin.Password)
	return nil
}
