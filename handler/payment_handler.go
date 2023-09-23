package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPaymentHealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Welcome to ev-payment-service"))
	return
}

func CreatePaymentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Payment created"))
	return
}

func GetPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Payment retrieved"))
	return
}

func UpdatePaymentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Payment updated"))
	return
}

func CreateResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}
