package handler

import (
	"ev-payment-service/config"
	"ev-payment-service/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// @Summary		Health Check
// @Description 	perform health check status
// @Tags 			Health Check
// @Accept 		json
// @Produce 		json
// @Success 		200	{object}	map[string]interface{}	"returns a welcome message"
// @Router			/payment/home [get]
func GetPaymentHealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Welcome to ev-payment-service"))
	return
}

func CreateResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

// @Summary		Get Invoice by user or provider (not sure pending stripe integration to know what is needed)
// @Description	get Invoice by user or provider
// @Tags			Invoice
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.Invoice
// @Router			/invoice [get]
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

// @Summary		Create Invoice by user or provider (not sure pending stripe integration to know what is needed)
// @Description	create Invoice by user or provider
// @Tags			Invoice
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.Invoice
// @Router			/invoice [post]
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
