package handler

import (
	"encoding/json"
	"ev-payment-service/config"
	"ev-payment-service/controller/user"
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

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	bookingId, err := strconv.Atoi(c.Param("booking_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("booking id must be an integer"))
	}

	userPayments, err := userpayment.UserControllerObj.GetUserPaymentInfo(uint(bookingId))
	if err != nil {
		logrus.WithField("err", err).Error("error getting user payment")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, userPayments)
	return
}

func CreateUserPaymentHandler(c *gin.Context) {
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

	if err := json.NewDecoder(c.Request.Body).Decode(&userPayment); err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	userPayment.UserEmail = user.Email

	stripe, err := userpayment.UserControllerObj.CreateUserPayment(&userPayment, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if stripe != "" {
		c.JSON(http.StatusOK, gin.H{"stripe": stripe, "userPayment": userPayment})
		return
	}

}

func UpdateUserPaymentHandler(c *gin.Context) {
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

	err = c.BindJSON(&userPayment)
	if err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	userPayment.UserEmail = user.Email

	err = userpayment.UserControllerObj.UpdateUserPayment(userPayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, CreateResponse("success"))
	return
}

func DeleteUserPaymentHandler(c *gin.Context) {
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
	userPayment.BookingId = uint(bookingId)

	err = userpayment.UserControllerObj.DeleteUserPayment(uint(bookingId))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, CreateResponse("success"))
	return
}

func CompleteUserPaymentHandler(c *gin.Context) {
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
	userPayment.BookingId = uint(bookingId)

	err = userpayment.UserControllerObj.CompleteUserPayment(&userPayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	//c.JSON(http.StatusOK, gin.H{"userPayment": userPayment})
	c.JSON(http.StatusOK, CreateResponse("success"))
	return
}
