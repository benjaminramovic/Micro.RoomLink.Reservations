package models

import (
	"time"
)

type Reservation struct {
	Id string `json:"id" bson:"_id,omitempty"`
	GuestId string `json:"guest_id" bson:"guest_id"`
	RoomId string `json:"room_id" bson:"room_id"`
	CheckIn time.Time `json:"check_in" bson:"check_in"`
	CheckOut time.Time `json:"check_out" bson:"check_out"`
	TotalPrice float64 `json:"total_price" bson:"total_price"`
	Status string `json:"status" bson:"status"`
}