package dto

import "fmt"

type ProviderPayment struct {
	Payment
	ProviderBillingPeriod
	TotalCommission float64 `bson:"TotalCommission,omitempty" json:"TotalCommission,omitempty"`
	CommissionRate  float64 `bson:"CommissionRate,omitempty" json:"CommissionRate,omitempty"`
	ProviderId      uint    `bson:"ProviderId,omitempty" json:"ProviderId,omitempty"`
	UserEmail       string  `bson:"UserEmail,omitempty" json:"UserEmail,omitempty"`
	PaymentStatus   string  `bson:"PaymentStatus,omitempty" json:"PaymentStatus,omitempty"`
}

func (p ProviderPayment) SetCompleteStatus() error {
	if p.PaymentStatus == "waiting" {
		p.PaymentStatus = "completed"
		return nil
	} else {
		return fmt.Errorf("provider payment is not complete status")
	}
}

type ProviderBillingPeriod struct {
	BillingMonth uint `bson:"column:billingMonth,omitempty"`
	BillingYear  uint `bson:"column:billingYear,omitempty"`
}
