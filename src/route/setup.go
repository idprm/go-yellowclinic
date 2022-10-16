package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/controller"
)

func Setup(app *fiber.App) {

	v1 := app.Group("v1")
	v1.Post("auth", controller.AuthHandler)

}
