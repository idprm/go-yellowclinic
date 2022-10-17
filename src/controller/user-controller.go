package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/model"
)

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var user model.User
	user.ID = uint64(id)
	database.Datasource.DB().First(&user)

	return c.Status(fiber.StatusOK).JSON(user)
}
