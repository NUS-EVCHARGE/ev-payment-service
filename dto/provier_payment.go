package dto

type ProviderPayment struct {
	Payment
	ProviderBillingPeriod
	TotalCommission float64 `bson:"TotalCommission,omitempty" json:"TotalCommission,omitempty"`
	CommissionRate  float64 `bson:"CommissionRate,omitempty" json:"CommissionRate,omitempty"`
	ProviderId      uint    `bson:"ProviderId,omitempty" json:"ProviderId,omitempty"`
	UserEmail       string  `bson:"UserEmail,omitempty" json:"UserEmail,omitempty"`
	Status          string  `bson:"Status,omitempty" json:"Status,omitempty"`
}

type ProviderBillingPeriod struct {
	BillingMonth uint `bson:"column:billingMonth,omitempty"`
	BillingYear  uint `bson:"column:billingYear,omitempty"`
}
