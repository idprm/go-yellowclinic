package controller

import (
	"strconv"

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
	id, _ := strconv.Atoi(c.Params("id"))

	var doctor model.Doctor
	doctor.ID = uint(id)
	database.Datasource.DB().First(&doctor)

	return c.Status(fiber.StatusOK).JSON(doctor)
}
