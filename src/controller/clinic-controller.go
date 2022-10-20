package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/model"
)

func GetAllClinic(c *fiber.Ctx) error {

	var clinics []model.Clinic
	database.Datasource.DB().Where("is_active", true).Order("id desc").Find(&clinics)

	return c.Status(fiber.StatusOK).JSON(clinics)
}

func GetClinic(c *fiber.Ctx) error {
	username := c.Params("username")

	var clinic model.Clinic
	database.Datasource.DB().Where("username", username).First(&clinic)

	return c.Status(fiber.StatusOK).JSON(clinic)
}
