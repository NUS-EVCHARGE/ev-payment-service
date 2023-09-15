package dto

type ProviderPayment struct {
	Payment
	TotalCommission float64
	ComissionRate   float64
	ProviderId      uint
}
