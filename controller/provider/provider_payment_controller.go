package controller

import "ev-payment-service/dto"

type ProviderPaymentController interface {
	GetProviderPaymentInfo(providerId uint) (dto.ProviderPayment, error)
	CreateProviderPayment(providerPayment dto.ProviderPayment) error
}
