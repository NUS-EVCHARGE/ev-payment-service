package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TotalBill float64            `bson:"totalBill,omitempty"`
	FinalBill float64            `bson:"finalBill,omitempty"`
}
