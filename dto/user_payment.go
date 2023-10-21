package dto

type UserPayment struct {
	Payment
	Booking
	ChargerAddress string  `bson:"chargerAddress,omitempty" json:"chargerAddress,omitempty"`
	RateID         uint    `bson:"rateId,omitempty" json:"rateId,omitempty"`
	NormalRate     float64 `bson:"normalRate,omitempty" json:"normalRate,omitempty"`
	BookingId      uint    `bson:"bookingId,omitempty" json:"bookingId,omitempty"`
	UserEmail      string  `bson:"userEmail,omitempty" json:"userEmail,omitempty"`
	PaymentStatus  string  `bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
}

func (u *UserPayment) SetCompleteStatus() {
	u.PaymentStatus = "completed"
}
