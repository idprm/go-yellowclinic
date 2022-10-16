package controller

import "github.com/gofiber/fiber/v2"

func FrontHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "messsage": "Welcome to yellowclinic"})
}

func AuthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Active"})
}
