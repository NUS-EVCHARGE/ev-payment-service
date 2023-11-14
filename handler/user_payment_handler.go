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

// @Summary Get User Payment
// @Description get user payment
// @Tags user payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Param booking_id path int true "booking id"
// @Success 200 {object} []dto.UserPayment "returns a user payment object"
// @Router /payment/user/booking/{booking_id} [get]
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

// @Summary Create User Payment
// @Description create user payment
// @Tags user payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Router /payment/user [post]
// @Success 200 {object} dto.UserPayment "returns a stripe id"
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

// @Summary Update User Payment
// @Description update user payment
// @Tags user payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Router /payment/user/booking/{booking_id} [put]
// @Success 200 {object} string "returns a success message"
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

// @Summary Delete User Payment
// @Description delete user payment
// @Tags user payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Router /payment/user/booking/{booking_id} [delete]
// @Success 200 {object} string "returns a success message"
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

// @Summary Complete User Payment and save a record as invoice
// @Description complete user payment
// @Tags user payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Router /payment/user/completed [Post]
// @Success 200 {object} string "returns a success message"
func CompleteUserPaymentHandler(c *gin.Context) {
	var (
		userPayment dto.UserPayment
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)

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

	err = userpayment.UserControllerObj.CompleteUserPayment(&userPayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	//c.JSON(http.StatusOK, gin.H{"userPayment": userPayment})
	c.JSON(http.StatusOK, CreateResponse("success"))
	return
}

// @Summary Get All User Payments by user
// @Description get all user payments by user email address
// @Tags user payment
// @Accept json
// @Produce json
// @Router /payment/user/booking [get]
// @Param authentication header string true "jwtToken of the user"
// @Success 200 {object} map[string][]dto.UserPayment "returns a map of user payment for, pending and completed"
func GetAllUserPaymentHandler(c *gin.Context) {
	var (
		user userDto.User
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	user, err := helper.GetUser(config.GetUserUrl, tokenStr)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	userPayment, err := userpayment.UserControllerObj.GetAllUserPayments(tokenStr, user)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user payment")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, userPayment)
	return

}
