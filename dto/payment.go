package dto

type Payment struct {
	//ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TotalBill float64 `bson:"totalBill,omitempty"`
	FinalBill float64 `bson:"finalBill,omitempty"`
}
