package mongoo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SellerID           primitive.ObjectID `bson:"seller_id" json:"seller_id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	Price              float64            `bson:"price" json:"price"`
	Quantity           int                `bson:"quantity" json:"quantity"`
	Wide               float64            `bson:"wide" json:"wide"`
	Height             float64            `bson:"height" json:"height"`
	Length             float64            `bson:"length" json:"length"`
	Weight             float64            `bson:"weight" json:"weight"`
	Category           string             `bson:"category" json:"category"`
	ChildCategory      string             `bson:"child_category" json:"child_category"`
	GrandchildCategory string             `bson:"grandchild_category" json:"grandchild_category"`
	ProfilePictureURL  []string           `json:"profile_picture_url,omitempty" bson:"profile_picture_url,omitempty"`
}
