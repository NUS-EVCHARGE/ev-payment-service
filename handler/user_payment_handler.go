package handler

import (
	"ev-payment-service/config"
	usercontroller "ev-payment-service/controller/user"
	"ev-payment-service/dto"
	"ev-payment-service/helper"
	"fmt"
	userDto "github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetUserPaymentHandler(c *gin.Context) {
	var (
		user        userDto.User
		userPayment dto.UserPayment
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	user, err := helper.GetUser(config.GetUserUrl, tokenStr)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	bookingId, err := strconv.Atoi(c.Param("booking_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("booking id must be an integer"))
	}

	userPayment.UserEmail = user.Email
	userPayment, err = usercontroller.UserPaymentControllerObj.GetUserPaymentInfo(uint(bookingId))
	if err != nil {
		logrus.WithField("err", err).Error("error getting user payment")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, userPayment)
	return
}
