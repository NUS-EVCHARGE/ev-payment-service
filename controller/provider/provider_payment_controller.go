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
}

var (
	ProviderPaymentControllerObj ProviderPaymentController
)

type ProviderPaymentControllerImpl struct {
	stripeKey string
}

func (p ProviderPaymentControllerImpl) GetProviderPaymentInfo(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error) {
	providerPayment, err := dao.Db.GetProviderPaymentEntry(providerId, billingPeriod)
	if err != nil {
		return nil, err
	}
	return providerPayment, nil
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

	if providerPayment.Status != "waiting" {
		return fmt.Errorf("provider payment has already been completed")
	}

	providerPayment.Status = "completed"
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

func NewProviderPaymentController(stripeKey string) {
	ProviderPaymentControllerObj = &ProviderPaymentControllerImpl{
		stripeKey: stripeKey,
	}
}
