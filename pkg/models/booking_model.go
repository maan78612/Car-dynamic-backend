package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bookings struct {
	// _id is bson format So it takes ID  in json response instead of _id
	ID               primitive.ObjectID `bson:"_id"`
	User_id          *string            `json:"user_id"`
	Created_at       time.Time          `json:"created_at"`
	Duration_in_days int                `json:"duration_in_days" validate:"required"`
	Start_time       time.Time          `json:"start_time" validate:"required"`
	End_time         time.Time          `json:"end_time"`
}
