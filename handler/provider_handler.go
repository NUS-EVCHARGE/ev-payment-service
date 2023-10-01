package handler

import (
	"encoding/json"
	"ev-payment-service/config"
	provider "ev-payment-service/controller/provider"
	"ev-payment-service/dto"
	"ev-payment-service/helper"
	"fmt"
	userDto "github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateProviderPaymentHandler(c *gin.Context) {
	var (
		user            userDto.User
		providerPayment dto.ProviderPayment
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	user, err := helper.GetUser(config.GetUserUrl, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&providerPayment); err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Error decoding request body"))
		return
	}

	providerPayment.UserEmail = user.Email

	stripe, err := provider.ProviderPaymentControllerObj.CreateProviderPayment(&providerPayment, tokenStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if stripe != "" {
		c.JSON(http.StatusOK, gin.H{"stripe": stripe, "provider_payment": providerPayment})
		return
	}
}

func GetProviderPaymentHandler(c *gin.Context) {

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Error getting user"))
		return
	}

	providerId, err := strconv.Atoi(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Provider id must be an integer"))
		return
	}

	billingMonth, _ := c.GetQuery("billing_month")
	billingYear, _ := c.GetQuery("billing_year")
	if billingMonth == "" || billingYear == "" {
		c.JSON(http.StatusBadRequest, CreateResponse("Billing month and year must be provided"))
		return
	}

	billingMonthInt, err := strconv.Atoi(billingMonth)
	billingYearInt, err := strconv.Atoi(billingYear)
	if billingMonthInt < 1 || billingMonthInt > 12 {
		c.JSON(http.StatusBadRequest, CreateResponse("Billing month must be between 1 and 12"))
		return
	}

	billingPeriod := dto.ProviderBillingPeriod{
		BillingMonth: uint(billingMonthInt),
		BillingYear:  uint(billingYearInt),
	}

	providerPayment, err := provider.ProviderPaymentControllerObj.GetProviderPaymentInfo(uint(providerId), billingPeriod)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, providerPayment)
	return
}

func UpdateProviderPaymentHandler(c *gin.Context) {

	var (
		user            userDto.User
		providerPayment dto.ProviderPayment
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	user, err := helper.GetUser(config.GetUserUrl, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&providerPayment); err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Error decoding request body"))
		return
	}

	providerPayment.UserEmail = user.Email

	err = provider.ProviderPaymentControllerObj.UpdateProviderPayment(providerPayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, CreateResponse("Payment updated"))
	return
}

func DeleteProviderPaymentHandler(c *gin.Context) {

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	_, err := helper.GetUser(config.GetUserUrl, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Error getting user"))
		return
	}

	providerId, err := strconv.Atoi(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Provider id must be an integer"))
		return
	}

	billingMonth, _ := c.GetQuery("billing_month")
	billingYear, _ := c.GetQuery("billing_year")
	if billingMonth == "" || billingYear == "" {
		c.JSON(http.StatusBadRequest, CreateResponse("Billing month and year must be provided"))
		return
	}

	billingMonthInt, err := strconv.Atoi(billingMonth)
	billingYearInt, err := strconv.Atoi(billingYear)
	if billingMonthInt < 1 || billingMonthInt > 12 {
		c.JSON(http.StatusBadRequest, CreateResponse("Billing month must be between 1 and 12"))
		return
	}

	billingPeriod := dto.ProviderBillingPeriod{
		BillingMonth: uint(billingMonthInt),
		BillingYear:  uint(billingYearInt),
	}

	err = provider.ProviderPaymentControllerObj.DeleteProviderPayment(uint(providerId), billingPeriod)

	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, CreateResponse("Payment deleted"))
	return
}

func CompleteProviderPaymentHandler(c *gin.Context) {

	var (
		user            userDto.User
		providerPayment dto.ProviderPayment
	)

	tokenStr := c.GetHeader("Authentication")

	// Get User information
	user, err := helper.GetUser(config.GetUserUrl, tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	providerId, err := strconv.Atoi(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Provider id must be an integer"))
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&providerPayment); err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse("Error decoding request body"))
		return
	}

	providerPayment.UserEmail = user.Email
	providerPayment.ProviderId = uint(providerId)

	err = provider.ProviderPaymentControllerObj.CompleteProviderPayment(&providerPayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, CreateResponse("Payment completed"))
	return
}
