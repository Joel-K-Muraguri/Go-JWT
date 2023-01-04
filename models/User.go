package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{

	ID  primitive.ObjectID   `bson:" id"`
	First_Name *string       `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name *string        `json:"last_name" validate:"required,min=2,max=100"`
	Phone_Number *string  `json:"phone" validate:"required"`
	User_ID string    `json:"user_id"`
	User_type *string  `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Email *string       `json:"email" validate:"email,required"`
	Password *string    `json:"Password" validate:"required,min=6"`
	Token *string    `json:"token"`
	Refresh_token  *string `json:"refresh_token"`
	Created_at  time.Timer  `json:"created_at"`
	Updated_at  time.Timer  `json:"updated_at"`

}