package dto

import (
	"fmt"
	"gorm.io/gorm"
)

type ProviderPayment struct {
	gorm.Model
	Payment
	ProviderBillingPeriod
	TotalCommission float64 `gorm:"column:total_commission" bson:"TotalCommission,omitempty" json:"TotalCommission,omitempty"`
	CommissionRate  float64 `gorm:"column:commission_rate" bson:"CommissionRate,omitempty" json:"CommissionRate,omitempty"`
	ProviderId      uint    `gorm:"column:provider_id" bson:"ProviderId,omitempty" json:"ProviderId,omitempty"`
	UserEmail       string  `gorm:"column:user_email" bson:"UserEmail,omitempty" json:"UserEmail,omitempty"`
	PaymentStatus   string  `gorm:"column:payment_status" bson:"PaymentStatus,omitempty" json:"PaymentStatus,omitempty"`
}

type ProviderEarnings struct {
	TotalEarnings   float64 `gorm:"column:total_earnings" bson:"TotalEarnings,omitempty" json:"TotalEarnings,omitempty"`
	TotalCommission float64 `gorm:"column:total_commission" bson:"TotalCommission,omitempty" json:"TotalCommission,omitempty"`
	NetEarnings     float64 `gorm:"column:net_earnings" bson:"NetEarnings,omitempty" json:"NetEarnings,omitempty"`
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
	BillingMonth uint `gorm:"column:billing_month" bson:"column:billingMonth,omitempty"`
	BillingYear  uint `gorm:"column:billing_year " bson:"column:billingYear,omitempty"`
}

func (ProviderPayment) TableName() string {
	return "provider_payment_tab"
}
