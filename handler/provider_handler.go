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

// @Summary Create Provider Payment
// @Description create provider payment
// @Tags provider payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Param provider_id path int true "provider id"
// @Success 200 {object} dto.ProviderPayment "returns a provider payment object with a stripe key"
// @Router /payment/provider [get]
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

// @Summary Get Provider Payment
// @Description get provider payment
// @Tags provider payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Param provider_id path int true "provider id"
// @Param billing_month query int true "billing month"
// @Param billing_year query int true "billing year"
// @Success 200 {object} dto.ProviderPayment
// @Router /payment/provider/{provider_id} [get]
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

// @Summary Update Provider Payment
// @Description update provider payment
// @Tags provider payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Param provider_id path int true "provider id"
// @success 200 {object} string "returns a string"
// @Router /payment/provider/{provider_id} [put]
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

// @Summary Delete Provider Payment
// @Description delete provider payment
// @Tags provider payment
// @Accept json
// @Produce json
// @Param authentication header string true "jwtToken of the user"
// @Param provider_id path int true "provider id"
// @Param billing_month query int true "billing month"
// @Param billing_year query int true "billing year"
// @success 200 {object} string "returns a string"
// @Router /payment/provider/{provider_id} [delete]
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

// @summary Complete Provider Payment
// @description complete provider payment
// @tags provider payment
// @accept json
// @produce json
// @param authentication header string true "jwtToken of the user"
// @param provider_id path int true "provider id"
// @success 200 {object} string "returns a string"
// @router /payment/provider/completed/{provider_id} [put]
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

func GetProviderTotalEarningsByProviderIDHandler(c *gin.Context) {
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

	totalEarnings, err := provider.ProviderPaymentControllerObj.GetProviderTotalEarningsByProviderID(uint(providerId))
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateResponse(fmt.Sprintf("%v", err)))
		return
	}

	c.JSON(http.StatusOK, totalEarnings)
	return
}
