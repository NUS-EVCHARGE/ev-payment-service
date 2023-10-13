package dto

import "time"

type Booking struct {
	ID        uint      `bson:"column:id,omitempty" json:"id,omitempty"`
	ChargerId uint      `bson:"column:charger_id,omitempty" json:"charger_id,omitempty"`
	Email     string    `bson:"column:email,omitempty" json:"email,omitempty"`
	StartTime time.Time `bson:"column:start_time,omitempty" json:"start_time,omitempty"`
	EndTime   time.Time `bson:"column:end_time,omitempty" json:"end_time,omitempty"`
	Status    string    `bson:"column:status,omitempty" json:"status,omitempty"`
}
