package user

import (
	"amanahkirim/db/mongoo"
	"amanahkirim/graphql/user/utils"
	rootUtils "amanahkirim/graphql/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	username := p.Args["username"].(string)
	password := p.Args["password"].(string)
	confirmPassword := p.Args["confirm_password"].(string)
	collection := mongoo.ClientUser.Database("userdb").Collection("users")

	if password != confirmPassword {
		return nil, errors.New("Password doesn't match")
	}

	_, err := IsUserExist(username)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("Failed to encrypt password")
	}

	user := mongoo.User{
		Username: username,
		Password: hashPassword,
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, errors.New("Failed to create new user")
	}

	jwtData := rootUtils.JwtData{
		ID: result.InsertedID.(primitive.ObjectID).Hex(),
	}

	token, err := utils.GenerateJWT(jwtData)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Failed to generate token, please re-login")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func Login(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username, ok := p.Args["username"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'username' parameter")
	}

	password, ok := p.Args["password"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'password' parameter")
	}
	collection := mongoo.ClientUser.Database("userdb").Collection("users")
	filter := bson.D{{Key: "username", Value: username}}

	var result mongoo.User

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("username/password are incorect")
	} else if err != nil {
		return false, err
	}

	err = utils.ComparePassword([]byte(result.Password), password)

	if err != nil {
		return nil, errors.New("username/password are incorect")
	}

	jwtData := rootUtils.JwtData{
		ID: result.ID.Hex(),
	}

	token, err := utils.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("Failed to login")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func NewBuyer(p graphql.ResolveParams) (interface{}, error) {
	userCtx, err := rootUtils.GetUserFromContext(p)
	if err != nil {
		return nil, err
	}
	userID := userCtx.ID
	sellerID := userCtx.SellerID

	BuyerData, _ := GetBuyerByUserId(userID)

	if BuyerData != nil {
		return nil, errors.New("You are already in existing Buyer")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addressID, _ := createAddress(ctx, p)

	buyerID, err := createBuyer(ctx, p, userID, addressID)
	if err != nil {
		return nil, err
	}

	err = updateUserFieldById(ctx, userID, "buyer_id", buyerID)
	if err != nil {
		return nil, err
	}

	buyerIDStr := buyerID.Hex()

	jwtData := rootUtils.JwtData{
		ID:       userID,
		BuyerID:  &buyerIDStr,
		SellerID: sellerID,
	}

	token, err := utils.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("Failed to generate JWT token")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func IsUserExist(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoo.ClientUser.Database("userdb").Collection("users")
	filter := bson.D{{Key: "username", Value: username}}

	var result mongoo.User

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, errors.New("Username already exist")
	}
}

func GetUserById(id string) (*mongoo.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoo.ClientUser.Database("userdb").Collection("users")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}

	var result mongoo.User

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("User not found")
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}
func GetBuyerByUserId(id string) (*mongoo.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoo.ClientUser.Database("userdb").Collection("buyers")
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

func createBuyer(ctx context.Context, p graphql.ResolveParams, userId string, addressID *primitive.ObjectID) (*primitive.ObjectID, error) {
	collectionBuyer := mongoo.ClientUser.Database("userdb").Collection("buyers")

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, errors.New("name is required")
	}
	phoneNumbers, ok := p.Args["phone_numbers"].(string)
	if !ok {
		return nil, errors.New("phone_numbers is required")
	}
	dateOfBirth, ok := p.Args["date_of_birth"].(string)
	if !ok {
		return nil, errors.New("date_of_birth is required")
	}
	profilePictureURL, ok := p.Args["profile_picture_url"].(string)

	userIdObject, err := primitive.ObjectIDFromHex(userId)

	layout := "02/01/2006"
	parsedDateOfBirth, err := time.Parse(layout, dateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date_of_birth format, use dd/mm/yyyy")
	}

	buyer := mongoo.Buyer{
		ID:                primitive.NewObjectID(),
		Name:              name,
		PhoneNumbers:      phoneNumbers,
		DateOfBirth:       parsedDateOfBirth,
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

	updateResult, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("Failed to update user profile")
	}

	fmt.Printf("Matched %v document(s) and modified %v document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}
