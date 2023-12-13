package dto

import "gorm.io/gorm"

type UserPayment struct {
	gorm.Model
	Payment
	Booking
	ChargerAddress string  `gorm:"column:charger_address" bson:"chargerAddress,omitempty" json:"chargerAddress,omitempty"`
	RateID         uint    `gorm:"column:rate_id" bson:"rateId,omitempty" json:"rateId,omitempty"`
	NormalRate     float64 `gorm:"column:normal_rate" bson:"normalRate,omitempty" json:"normalRate,omitempty"`
	ProviderId     uint    `gorm:"column:provider_id" bson:"providerId,omitempty" json:"providerId,omitempty"`
	BookingId      uint    `gorm:"column:booking_id" bson:"bookingId,omitempty" json:"bookingId,omitempty"`
	UserEmail      string  `gorm:"column:user_email" bson:"userEmail,omitempty" json:"userEmail,omitempty"`
	PaymentStatus  string  `gorm:"column:payment_status" bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
}

func (UserPayment) TableName() string {
	return "user_payment_tab"
}

func (u *UserPayment) SetCompleteStatus() {
	u.PaymentStatus = "completed"
}
