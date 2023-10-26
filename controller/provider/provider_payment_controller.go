package controller

import (
	"ev-payment-service/config"
	"ev-payment-service/dao"
	"ev-payment-service/dto"
	helper "ev-payment-service/helper"
	"fmt"
)

type ProviderPaymentController interface {
	GetProviderPaymentInfo(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error)
	CreateProviderPayment(providerPayment *dto.ProviderPayment, token string) (string, error)
	UpdateProviderPayment(providerPayment dto.ProviderPayment) error
	DeleteProviderPayment(id uint, billingPeriod dto.ProviderBillingPeriod) error
	CompleteProviderPayment(providerPayment *dto.ProviderPayment) error
	GetProviderTotalEarningsByProviderID(providerId uint) (dto.ProviderEarnings, error)
}

var (
	ProviderPaymentControllerObj ProviderPaymentController
)

type ProviderPaymentControllerImpl struct {
	stripeKey string
}

func (p ProviderPaymentControllerImpl) GetProviderPaymentInfo(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error) {
	return dao.Db.GetProviderPaymentEntry(providerId, billingPeriod)
}

func (p ProviderPaymentControllerImpl) UpdateProviderPayment(providerPayment dto.ProviderPayment) error {
	return dao.Db.UpdateProviderPaymentEntry(&providerPayment)
}

func (p ProviderPaymentControllerImpl) DeleteProviderPayment(id uint, billingPeriod dto.ProviderBillingPeriod) error {
	return dao.Db.DeleteProviderPaymentEntry(id, billingPeriod)
}

func (p ProviderPaymentControllerImpl) CompleteProviderPayment(providerPayment *dto.ProviderPayment) error {
	it, err := p.GetProviderPaymentInfo(providerPayment.ProviderId, providerPayment.ProviderBillingPeriod)
	if err != nil {
		return err
	}

	if len(it) > 0 {
		providerPayment = &it[0]
	} else {
		return fmt.Errorf("provider payment not found")
	}

	if err := providerPayment.SetCompleteStatus(); err != nil {
		return err
	}

	return dao.Db.UpdateProviderPaymentEntry(providerPayment)
}

func (p ProviderPaymentControllerImpl) CreateProviderPayment(providerPayment *dto.ProviderPayment, token string) (string, error) {

	providerResult, err := helper.GetProvider(config.GetProviderUrl, token)
	if providerResult.UserEmail != providerPayment.UserEmail {
		return "", fmt.Errorf("user is not the owner of the provider")
	}

	// Check if there is no duplication in storage
	providerPaymentStored, err := p.GetProviderPaymentInfo(providerPayment.ProviderId, providerPayment.ProviderBillingPeriod)
	if providerPaymentStored != nil && len(providerPaymentStored) > 0 {
		return "", fmt.Errorf("provider payment already exist")
	}

	if providerPayment.Coupon != "" {
		providerPayment.CommissionRate = 0.05
	} else {
		providerPayment.CommissionRate = 0.1
	}

	providerPayment.TotalCommission = providerPayment.TotalBill * providerPayment.CommissionRate
	providerPayment.FinalBill = providerPayment.TotalBill - providerPayment.TotalCommission
	providerPayment.PaymentStatus = "pending"

	stripeClientSecret, err := helper.CreateStripeSecret(providerPayment.FinalBill)
	if err != nil {
		return "", err
	}

	dbErr := dao.Db.CreateProviderPaymentEntry(providerPayment)
	if dbErr != nil {
		return "", dbErr
	} else {
		return stripeClientSecret, nil
	}
}

func (p ProviderPaymentControllerImpl) GetProviderTotalEarningsByProviderID(providerId uint) (dto.ProviderEarnings, error) {
	providerPaymentResult, err := dao.Db.GetUserPaymentByProviderId(providerId)
	if err != nil {
		return dto.ProviderEarnings{
			TotalEarnings:   0,
			TotalCommission: 0,
			NetEarnings:     0,
		}, fmt.Errorf("error getting provider payment: %v", err)
	}

	var totalEarnings float64
	var totalCommission float64
	for _, providerPayment := range providerPaymentResult {
		totalEarnings = totalEarnings + providerPayment.TotalBill
		totalCommission = totalCommission + (providerPayment.TotalBill * 0.05)
	}

	return dto.ProviderEarnings{
		TotalEarnings:   totalEarnings,
		TotalCommission: totalCommission,
		NetEarnings:     totalEarnings - totalCommission,
	}, nil
}

func NewProviderPaymentController(stripeKey string) {
	ProviderPaymentControllerObj = &ProviderPaymentControllerImpl{
		stripeKey: stripeKey,
	}
}
