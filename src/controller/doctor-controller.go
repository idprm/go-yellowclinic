package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/model"
)

func GetAllDoctor(c *fiber.Ctx) error {

	var doctors []model.Doctor
	database.Datasource.DB().Order("end desc").Find(&doctors)

	return c.Status(fiber.StatusOK).JSON(doctors)
}

func GetDoctor(c *fiber.Ctx) error {
	username := c.Params("username")

	var doctor model.Doctor
	database.Datasource.DB().Where("user_id", username).First(&doctor)

	return c.Status(fiber.StatusOK).JSON(doctor)
}
