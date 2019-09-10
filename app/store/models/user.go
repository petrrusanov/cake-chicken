package models

import (
	"time"
)

//Gender user gender
type Gender string

//Gender gender enum
const (
	Male   Gender = "male"
	Female Gender = "female"
)

//User user struct
type User struct {
	ID           string       `json:"id" bson:"_id"`
	Gender       Gender       `json:"gender"`
	DateOfBirth  time.Time    `json:"dateOfBirth,string" bson:"dateOfBirth"`
	Token        string       `json:"token"`
	CreatedAt    time.Time    `json:"createdAt,string,omitempty" bson:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt,string,omitempty" bson:"updatedAt"`
}
