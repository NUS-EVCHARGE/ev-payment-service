package dto

import "fmt"

type UserPayment struct {
	Payment
	Booking
	BookingId     uint   `bson:"bookingId,omitempty" json:"bookingId,omitempty"`
	UserEmail     string `bson:"userEmail,omitempty" json:"userEmail,omitempty"`
	PaymentStatus string `bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
}

func (u *UserPayment) SetCompleteStatus() error {
	if u.PaymentStatus == "pending" {
		u.PaymentStatus = "completed"
		return nil
	} else {
		return fmt.Errorf("user payment is not complete status")
	}
}
