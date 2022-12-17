package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/middleware"
	"github.com/idprm/go-yellowclinic/src/model"
)

type AuthRequest struct {
	Msisdn    string `query:"msisdn" validate:"required" json:"msisdn"`
	Name      string `query:"name" validate:"required" json:"name"`
	Voucher   string `query:"voucher" json:"voucher"`
	IpAddress string `query:"ip_address" json:"ip_address"`
}

type ErrorResponse struct {
	Field string
	Tag   string
	Value string
}

var validate = validator.New()

func ValidateAuth(req AuthRequest) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func FrontHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    false,
		"status":   fiber.StatusOK,
		"messsage": "Welcome to yellowclinic",
	})
}

func AuthHandler(c *fiber.Ctx) error {

	c.Accepts("application/json")

	req := new(AuthRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	errors := ValidateAuth(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": errors,
		})
	}

	var user model.User
	isExist := database.Datasource.DB().Where("msisdn", req.Msisdn).First(&user)

	if isExist.RowsAffected == 0 {
		database.Datasource.DB().Create(&model.User{
			Msisdn:        req.Msisdn,
			Name:          req.Name,
			LatestVoucher: req.Voucher,
			IpAddress:     req.IpAddress,
		})
	} else {
		user.Name = req.Name
		user.LatestVoucher = req.Voucher
		user.IpAddress = req.IpAddress
		database.Datasource.DB().Save(&user)
	}

	var usr model.User
	database.Datasource.DB().Where("msisdn", req.Msisdn).First(&usr)

	token, exp, err := middleware.GenerateJWTToken(usr)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"token": token,
		"exp":   exp,
		"user":  usr,
	})
}
