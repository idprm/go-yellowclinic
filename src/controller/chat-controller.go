package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/handler"
	"github.com/idprm/go-yellowclinic/src/model"
)

type ChatRequest struct {
	Msisdn     string `query:"msisdn" validate:"required" json:"msisdn"`
	ChannelUrl string `query:"channel_url" validate:"required" json:"channel_url"`
}

func ChatUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	req := new(ChatRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var user model.User
	database.Datasource.DB().Where("msisdn", req.Msisdn).First(&user)

	var chat model.Chat
	isChat := database.Datasource.DB().Where("user_id", user.ID).Preload("User").Preload("Doctor").First(&chat)

	if isChat.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&chat)
}

func ChatDoctor(c *fiber.Ctx) error {
	c.Accepts("application/json")

	channelUrl := c.Query("channel_url")

	var chat model.Chat
	database.Datasource.DB().Where("channel_url", channelUrl).Preload("User").Preload("Doctor").Preload("Order").First(&chat)

	return c.Status(fiber.StatusOK).JSON(&chat)
}

func ChatLeave(c *fiber.Ctx) error {

	c.Accepts("application/json")

	req := new(ChatRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var user model.User
	database.Datasource.DB().Where("msisdn", req.Msisdn).First(&user)

	var chat model.Chat
	database.Datasource.DB().Where("user_id", user.ID).First(&chat)

	leaveGroupChannel, isLeave, err := handler.SendbirdLeaveGroupChannel(chat.ChannelUrl, user.Msisdn)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   user.Msisdn,
		Action:   "LEAVE GROUP CHANNEL",
		Response: leaveGroupChannel,
	})

	if isLeave == true {
		//
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "Leaved"})
}

func ChatDelete(c *fiber.Ctx) error {
	c.Accepts("application/json")

	req := new(ChatRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var user model.User
	database.Datasource.DB().Where("msisdn", req.Msisdn).First(&user)

	var chat model.Chat
	database.Datasource.DB().Where("user_id", user.ID).Delete(&chat)

	deleteGroupChannel, err := handler.SendbirdDeleteGroupChannel(chat)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// insert to sendbird
	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   user.Msisdn,
		Action:   "DELETE GROUP CHANNEL",
		Response: deleteGroupChannel,
	})

	return c.Status(fiber.StatusOK).JSON(&chat)
}
