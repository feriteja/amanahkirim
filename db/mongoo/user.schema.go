package mongoo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" unique:"true"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email"`
	BuyerID  primitive.ObjectID `json:"buyer_id" bson:"buyer_id"`
	SellerID primitive.ObjectID `json:"seller_id" bson:"seller_id"`
}

type Buyer struct {
	ID                primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID            []primitive.ObjectID  `json:"user_id" bson:"user_id"`
	Name              string                `json:"name" bson:"name"`
	PhoneNumbers      string                `json:"phone_numbers" bson:"phone_numbers"`
	Address           []*primitive.ObjectID `json:"address,omitempty" bson:"address,omitempty"`
	DateOfBirth       time.Time             `json:"date_of_birth" bson:"date_of_birth"`
	ProfilePictureURL string                `json:"profile_picture_url" bson:"profile_picture_url"`
}

type Seller struct {
	ID                primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID            []primitive.ObjectID  `json:"user_id" bson:"user_id"`
	Name              string                `json:"name" bson:"name"`
	PhoneNumbers      string                `json:"phone_numbers" bson:"phone_numbers"`
	Address           []*primitive.ObjectID `json:"address,omitempty" bson:"address,omitempty"`
	ProfilePictureURL string                `json:"profile_picture_url" bson:"profile_picture_url"`
}

type Address struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Country    string             `json:"country" bson:"country"`
	City       string             `json:"city" bson:"city"`
	Province   string             `json:"province" bson:"province"`
	Regency    string             `json:"regency" bson:"regency"`
	PostalCode string             `json:"postal_code" bson:"postal_code"`
}
