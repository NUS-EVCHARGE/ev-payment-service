package dto

type UserPayment struct {
	Payment
	Booking
	BookingId uint   `bson:"bookingId,omitempty" json:"bookingId,omitempty"`
	UserEmail string `bson:"userEmail,omitempty" json:"userEmail,omitempty"`
}
