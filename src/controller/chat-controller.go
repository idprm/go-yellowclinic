package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/handler"
	"github.com/idprm/go-yellowclinic/src/model"
)

type ChatRequest struct {
	Msisdn     string `query:"msisdn" json:"msisdn"`
	Voucher    string `query:"voucher" json:"voucher"`
	ChannelUrl string `query:"channel_url" json:"channel_url"`
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

	var order model.Order
	isOrder := database.Datasource.DB().Where(&model.Order{UserID: user.ID, Voucher: req.Voucher}).First(&order)
	if isOrder.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(&order)
	}

	var chat model.Chat
	database.Datasource.DB().Joins("Order", database.Datasource.DB().Where(&model.Order{UserID: user.ID, Voucher: req.Voucher})).Where(&model.Chat{UserID: user.ID, IsLeave: false}).Preload("User").Preload("Doctor").Preload("Order").First(&chat)

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

	leaveGroupChannel, _, err := handler.SendbirdLeaveGroupChannel(chat.ChannelUrl, user.Msisdn)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   user.Msisdn,
		Action:   "USER LEAVE GROUP CHANNEL",
		Response: leaveGroupChannel,
	})

	chat.IsLeave = true
	chat.LeaveAt = time.Now()
	database.Datasource.DB().Save(&chat)

	callback, err := handler.CallbackVoucher(user.LatestVoucher)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	database.Datasource.DB().Create(&model.Callback{
		Msisdn:   user.Msisdn,
		Action:   user.LatestVoucher,
		Response: callback,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "Leaved"})
}

func ChatLeaveDoctor(c *fiber.Ctx) error {
	c.Accepts("application/json")

	req := new(ChatRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var chat model.Chat
	database.Datasource.DB().Where("channel_url", req.ChannelUrl).Preload("Doctor").Preload("User").First(&chat)

	leaveGroupChannel, _, err := handler.SendbirdLeaveGroupChannel(chat.ChannelUrl, chat.Doctor.Username)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	muteUserChannel, err := handler.SendbirdMuteUser(chat.ChannelUrl, chat.User)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	database.Datasource.DB().Create(&[]model.Sendbird{{
		Msisdn:   chat.Doctor.Phone,
		Action:   "DOCTOR LEAVE GROUP CHANNEL",
		Response: leaveGroupChannel,
	}, {
		Msisdn:   chat.User.Msisdn,
		Action:   "USER MUTE GROUP CHANNEL",
		Response: muteUserChannel,
	}})

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
