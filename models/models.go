package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Password     string             `json:"password" validate:"required,min=6"`
	Email        string             `json:"email" validate:"email,required"`
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	UserType     string             `bson:"user_type" json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken string             `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	UserId       string             `bson:"user_id"`
}
