package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-yellowclinic/src/config"
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/handler"
	"github.com/idprm/go-yellowclinic/src/model"
	"github.com/idprm/go-yellowclinic/src/util"
)

const (
	actionCheckUser     = "GET/CHECK USER"
	actionCreateUser    = "AUTO CREATE USER"
	actionCheckGroup    = "GET/CHECK GROUP CHANNEL"
	actionDeleteGroup   = "AUTO DELETE GROUP CHANNEL"
	actionCreateGroup   = "AUTO CREATE GROUP CHANNEL"
	actionCreateNotif   = "AUTO CREATE NOTIF TO DOCTOR"
	actionCreateMessage = "AUTO CREATE MESSAGE DOCTOR"
)

type OrderRequest struct {
	Msisdn   string `json:"msisdn"`
	Voucher  string `json:"voucher"`
	DoctorID int    `json:"doctor_id"`
}

func OrderChat(c *fiber.Ctx) error {
	c.Accepts("application/json")

	req := new(OrderRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var count int64
	database.Datasource.DB().Model(&model.Order{}).Count(&count)

	var confLimit model.Config
	database.Datasource.DB().Where("name", "LIMIT_VOUCHER").First(&confLimit)

	intLimit, _ := strconv.ParseInt(confLimit.Value, 0, 64)

	if count >= intLimit {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error":   true,
			"message": "Batas Voucher telah habis",
			"status":  fiber.StatusOK,
		})
	}

	var user model.User
	database.Datasource.DB().Where("msisdn", req.Msisdn).First(&user)

	var doctor model.Doctor
	database.Datasource.DB().Where("id", req.DoctorID).First(&doctor)

	/**
	 * Check Order
	 */
	var order model.Order
	resultOrder := database.Datasource.DB().
		Where("voucher", user.LatestVoucher).
		First(&order)

	finishUrl := config.ViperEnv("APP_HOST") + "/chat"

	if resultOrder.RowsAffected == 0 {
		/**
		 * INSERT TO ORDER
		 */
		database.Datasource.DB().Create(&model.Order{
			UserID:   user.ID,
			DoctorID: doctor.ID,
			Number:   "ORD-" + util.TimeStamp(),
			Voucher:  user.LatestVoucher,
			Total:    10000,
		})

		/**
		 * SENDBIRD PROCESS
		 */
		err := sendbirdProcess(user.ID, doctor.ID, user.LatestVoucher)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
				"status":  fiber.StatusBadGateway,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error":        false,
			"message":      "Created Successful",
			"redirect_url": finishUrl,
			"status":       fiber.StatusCreated,
		})
	} else {

		var chat model.Chat
		resultCount := database.Datasource.DB().Joins("Order", database.Datasource.DB().Where(&model.Order{Voucher: user.LatestVoucher})).Where("is_leave", false).First(&chat)

		if resultCount.RowsAffected > 0 {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"error":        false,
				"message":      "Already Chat",
				"redirect_url": finishUrl,
				"status":       fiber.StatusCreated,
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":   true,
				"message": "Voucher sudah terpakai",
				"status":  fiber.StatusOK,
			})
		}

	}
}

func sendbirdProcess(userId uint64, doctorId uint, latestVoucher string) error {

	var order model.Order
	database.Datasource.DB().
		Where("user_id", userId).
		Where("doctor_id", doctorId).
		Where("voucher", latestVoucher).
		Preload("User").Preload("Doctor").
		First(&order)
	/**
	 * Check User Sendbird
	 */
	getUser, isUser, err := handler.SendbirdGetUser(order.User)
	if err != nil {
		return errors.New(err.Error())
	}

	/**
	 * Add User Sendbird
	 */
	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   order.User.Msisdn,
		Action:   actionCheckUser,
		Response: getUser,
	})

	if isUser == true {

		// create user sendbird
		createUser, err := handler.SendbirdCreateUser(order.User)
		if err != nil {
			return errors.New(err.Error())
		}
		// insert to sendbird
		database.Datasource.DB().Create(&model.Sendbird{
			Msisdn:   order.User.Msisdn,
			Action:   actionCreateUser,
			Response: createUser,
		})
	}

	var chat model.Chat
	resultChat := database.Datasource.DB().Where("user_id", order.User.ID).First(&chat)

	getChannel, isChannel, err := handler.SendbirdGetGroupChannel(chat)
	if err != nil {
		return errors.New(err.Error())
	}

	// insert to sendbird
	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   order.User.Msisdn,
		Action:   actionCheckGroup,
		Response: getChannel,
	})

	// check db if exist
	if resultChat.RowsAffected > 0 {
		// check channel if exist
		if isChannel == false {
			// delete channel sendbird
			deleteGroupChannel, err := handler.SendbirdDeleteGroupChannel(chat)
			if err != nil {
				return errors.New(err.Error())
			}
			// delete chat
			database.Datasource.DB().Where("user_id", order.User.ID).Delete(&chat)
			// insert to sendbirds
			database.Datasource.DB().Create(&model.Sendbird{
				Msisdn:   order.User.Msisdn,
				Action:   actionDeleteGroup,
				Response: deleteGroupChannel,
			})
		}
	}

	// create group
	createGroup, name, url, err := handler.SendbirdCreateGroupChannel(order.Doctor, order.User)
	if err != nil {
		return errors.New(err.Error())
	}
	// insert to sendbird
	database.Datasource.DB().Create(&model.Sendbird{
		Msisdn:   order.User.Msisdn,
		Action:   actionCreateGroup,
		Response: createGroup,
	})

	if name != "" && url != "" {
		// insert to chat
		database.Datasource.DB().Create(&model.Chat{
			OrderID:     order.ID,
			DoctorID:    order.Doctor.ID,
			UserID:      order.User.ID,
			ChannelName: name,
			ChannelUrl:  url,
		})

		var conf model.Config
		database.Datasource.DB().Where("name", "NOTIF_MESSAGE_DOCTOR").First(&conf)

		urlWeb := config.ViperEnv("APP_HOST") + "/chat/" + url
		replaceMessage := strings.NewReplacer("@v1", order.Doctor.Name, "@v2", order.User.Name, "@v3", urlWeb)
		message := replaceMessage.Replace(conf.Value)

		// NOTIF MESSAGE TO DOCTOR
		zenzifaNotif, err := handler.ZenzivaSendSMS(order.Doctor.Phone, message)
		if err != nil {
			return errors.New(err.Error())
		}
		// insert to zenziva
		database.Datasource.DB().Create(&model.Zenziva{
			Msisdn:   order.User.Msisdn,
			Action:   actionCreateNotif,
			Response: zenzifaNotif,
		})

		// auto message to user
		autoMessageDoctor, err := handler.SendbirdAutoMessageDoctor(url, order.Doctor, order.User)
		if err != nil {
			return errors.New(err.Error())
		}

		// insert to sendbird
		database.Datasource.DB().Create(&model.Sendbird{
			Msisdn:   order.User.Msisdn,
			Action:   actionCreateMessage,
			Response: autoMessageDoctor,
		})
	}

	return nil
}
