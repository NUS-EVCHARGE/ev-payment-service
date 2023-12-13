package dto

type Payment struct {
	TotalBill float64 `gorm:"column:total_bill" bson:"totalBill,omitempty"`
	FinalBill float64 `gorm:"column:final_bill" bson:"finalBill,omitempty"`
	Coupon    string  `gorm:"column:coupon" bson:"coupon,omitempty"`
}
