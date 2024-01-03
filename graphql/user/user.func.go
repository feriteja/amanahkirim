package user

import (
	"amanahkirim/db/mongoo"
	"amanahkirim/graphql/user/utils"
	"context"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	jwtData := utils.JwtData{
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

	jwtData := utils.JwtData{
		ID: result.ID.Hex(),
	}

	token, err := utils.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("Failed to login")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func NewProfile(p graphql.ResolveParams) (interface{}, error) {
	userValue := p.Context.Value("user")
	userClaims, ok := userValue.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to assert the type of 'user'")
		return nil, errors.New("Failed to assert the type of 'user'")
	}

	id, ok := userClaims["data"].(map[string]interface{})["ID"].(string)
	if !ok {
		log.Println("Failed to get the 'ID' field")
		return nil, errors.New("Failed to get the 'ID' field")
	}

	user, err := GetUserById(id)
	if err != nil {
		return nil, errors.New("User not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionProfile := mongoo.ClientUser.Database("userdb").Collection("profile")
	collectionUser := mongoo.ClientUser.Database("userdb").Collection("users")

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'name' parameter")
	}
	phoneNumbers, ok := p.Args["phone_numbers"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'phone_numbers' parameter")
	}
	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'email' parameter")
	}

	city, ok := p.Args["city"].(string)
	province, ok := p.Args["province"].(string)
	country, ok := p.Args["country"].(string)
	postalCode, ok := p.Args["postal_code"].(string)
	dateOfBirth, ok := p.Args["date_of_birth"].(string)
	profilePictureURL, ok := p.Args["profile_picture_url"].(string)

	profile := mongoo.Profile{
		ID:                primitive.NewObjectID(),
		Name:              name,
		Email:             email,
		PhoneNumbers:      phoneNumbers,
		City:              city,
		Province:          province,
		Country:           country,
		PostalCode:        postalCode,
		DateOfBirth:       dateOfBirth,
		ProfilePictureURL: profilePictureURL,
	}
	result, err := collectionProfile.InsertOne(ctx, profile)
	if err != nil {
		return nil, errors.New("Failed to create new Profile")
	}

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "profile_id", Value: result.InsertedID.(primitive.ObjectID)}}}}

	_, err = collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.New("Failed to update user profile")
	}

	jwtData := utils.JwtData{
		ID:        id,
		ProfileID: result.InsertedID.(primitive.ObjectID).Hex(),
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
