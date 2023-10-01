package dto

type Invoice struct {
	UserEmail     string `bson:"column:user_email,omitempty" json:"user_email,omitempty"`
	InvoiceNumber string `bson:"column:invoice_number,omitempty" json:"invoice_number,omitempty"`
	PaymentDate   string `bson:"column:payment_date,omitempty" json:"payment_date,omitempty"`
	BookingID     uint   `bson:"column:booking_id,omitempty" json:"booking_id,omitempty"`
	ProviderID    uint   `bson:"column:provider_id,omitempty" json:"provider_id,omitempty"`
}
