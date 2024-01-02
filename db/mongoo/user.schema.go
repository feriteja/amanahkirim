package mongoo

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username" unique:"true"`
	Password string             `bson:"password"`
}
