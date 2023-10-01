package handler

import (
	"ev-payment-service/config"
	"ev-payment-service/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetPaymentHealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Welcome to ev-payment-service"))
	return
}

func CreateResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

func GetInvoiceHandler(c *gin.Context) {
	var (
	//user userDto.User
	//userPayment dto.UserPayment
	)

	tokenStr := c.GetHeader("Authentication")
	bookingId, bookingProvided := c.GetQuery("booking_id")
	providerId, providerProvided := c.GetQuery("provider_id")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if bookingProvided {
		c.JSON(http.StatusOK, gin.H{
			"booking_id": bookingId,
		})
	} else if providerProvided {
		c.JSON(http.StatusOK, gin.H{
			"provider_id": providerId,
		})
	} else {
		c.JSON(http.StatusBadRequest, CreateResponse("booking_id or provider_id must be provided"))
		return
	}

	return
}

func CreateInvoiceHandler(c *gin.Context) {
	var (
	//user userDto.User
	//userPayment dto.UserPayment
	)

	tokenStr := c.GetHeader("Authentication")
	bookingId, bookingProvided := c.GetQuery("booking_id")
	providerId, providerProvided := c.GetQuery("provider_id")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)

	if err != nil {
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if bookingProvided {
		c.JSON(http.StatusOK, gin.H{
			"booking_id": bookingId,
		})
	} else if providerProvided {
		c.JSON(http.StatusOK, gin.H{
			"provider_id": providerId,
		})
	} else {
		c.JSON(http.StatusBadRequest, CreateResponse("booking_id or provider_id must be provided"))
		return
	}

	return
}
