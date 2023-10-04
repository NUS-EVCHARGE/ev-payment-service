package dto

type Payment struct {
	TotalBill float64 `bson:"totalBill,omitempty"`
	FinalBill float64 `bson:"finalBill,omitempty"`
	Coupon    string  `bson:"coupon,omitempty"`
}
