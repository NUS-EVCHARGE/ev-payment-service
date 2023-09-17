package dto

type UserPayment struct {
	Payment
	BookingId uint   `bson:"bookingId,omitempty"`
	UserEmail string `bson:"userEmail,omitempty"`
}
