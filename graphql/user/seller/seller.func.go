package seller

import (
	"amanahkirim/db/mongoo"
	"amanahkirim/graphql/user/utils"
	rootUtils "amanahkirim/graphql/utils"
	"context"
	"errors"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewSeller(p graphql.ResolveParams) (interface{}, error) {
	userCtx, err := rootUtils.GetUserFromContext(p)
	if err != nil {
		return nil, err
	}
	userID := userCtx.ID
	buyerID := userCtx.BuyerID

	sellerData, _ := GetSellerByUserId(userID)

	if sellerData != nil {
		return nil, errors.New("You are already in existing Seller")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addressID, _ := createAddress(ctx, p)

	sellerID, err := createSeller(ctx, p, userID, addressID)
	if err != nil {
		return nil, err
	}

	err = updateUserFieldById(ctx, userID, "seller_id", sellerID)
	if err != nil {
		return nil, err
	}

	sellerIDStr := sellerID.Hex()

	jwtData := rootUtils.JwtData{
		ID:       userID,
		BuyerID:  buyerID,
		SellerID: &sellerIDStr,
	}

	token, err := utils.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("Failed to generate JWT token")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func GetSellerByUserId(id string) (*mongoo.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoo.ClientUser.Database("userdb").Collection("sellers")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "user_id", Value: objectID}}

	var result mongoo.User

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("User not found")
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

func createAddress(ctx context.Context, p graphql.ResolveParams) (*primitive.ObjectID, error) {
	collectionAddress := mongoo.ClientUser.Database("userdb").Collection("address")

	// Extracting input parameters
	country, ok := p.Args["country"].(string)
	if !ok {
		return nil, errors.New("Country is required")
	}

	city, ok := p.Args["city"].(string)
	if !ok {
		return nil, errors.New("City is required")
	}

	province, ok := p.Args["province"].(string)
	if !ok {
		return nil, errors.New("Province is required")
	}

	regency, ok := p.Args["regency"].(string)
	if !ok {
		return nil, errors.New("Regency is required")
	}

	postalCode, ok := p.Args["postal_code"].(string)
	if !ok {
		return nil, errors.New("Postal code is required")
	}

	// Create address request
	request := mongoo.Address{
		Country:    country,
		City:       city,
		Province:   province,
		Regency:    regency,
		PostalCode: postalCode,
	}

	// Insert the address into the database
	result, err := collectionAddress.InsertOne(ctx, request)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Failed to create new address")
	}

	// Convert the inserted ID to primitive.ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert inserted ID to ObjectID")
	}

	// Return the inserted ID
	return &insertedID, nil
}

func createSeller(ctx context.Context, p graphql.ResolveParams, userId string, addressID *primitive.ObjectID) (*primitive.ObjectID, error) {
	collectionBuyer := mongoo.ClientUser.Database("userdb").Collection("sellers")

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, errors.New("name is required")
	}
	phoneNumbers, ok := p.Args["phone_numbers"].(string)
	if !ok {
		return nil, errors.New("phone_numbers is required")
	}
	profilePictureURL, ok := p.Args["profile_picture_url"].(string)

	userIdObject, err := primitive.ObjectIDFromHex(userId)

	buyer := mongoo.Seller{
		ID:                primitive.NewObjectID(),
		Name:              name,
		PhoneNumbers:      phoneNumbers,
		ProfilePictureURL: profilePictureURL,
		Address:           []*primitive.ObjectID{addressID},
		UserID:            []primitive.ObjectID{userIdObject},
	}

	result, err := collectionBuyer.InsertOne(ctx, buyer)
	if err != nil {
		log.Println(err)
		return &primitive.NilObjectID, errors.New("Failed to create new buyer")
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert inserted ID to ObjectID")
	}

	return &insertedID, nil
}

func updateUserFieldById(ctx context.Context, userID string, fieldName string, value interface{}) error {
	collectionUser := mongoo.ClientUser.Database("userdb").Collection("users")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fieldName, Value: value}}}}

	_, err = collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("Failed to update user")
	}

	return nil
}
