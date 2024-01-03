package mongoo

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username" unique:"true"`
	Password  string             `bson:"password"`
	ProfileID string             `bson:"profile_id"`
}

type Profile struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Email             string             `bson:"email"`
	PhoneNumbers      string             `bson:"phone_numbers"`
	City              string             `bson:"city"`
	Province          string             `bson:"province"`
	Country           string             `bson:"country"`
	PostalCode        string             `bson:"postal_code"`
	DateOfBirth       string             `bson:"date_of_birth"`
	ProfilePictureURL string             `bson:"profile_picture_url"`
}
