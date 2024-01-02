package user

import (
	"amanahkirim/db/mongoo"
	"amanahkirim/graphql/user/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(p graphql.ResolveParams) (interface{}, error) {

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'name' parameter")
	}

	user, err := utils.GenerateJWT(name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	userID := p.Args["id"].(int)
	newName := p.Args["newName"].(string)

	user, ok := Users[userID]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", userID)
	}

	user.Username = newName
	return user, nil
}

func GetAllUsers(p graphql.ResolveParams) (interface{}, error) {
	return Users, nil
}

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

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, errors.New("Failed to create new user")
	}

	token, err := utils.GenerateJWT(username)
	if err != nil {
		return nil, errors.New("Failed to login")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func Login(p graphql.ResolveParams) (interface{}, error) {
	username, ok := p.Args["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'username' parameter")
	}
	// password, ok := p.Args["password"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid or missing 'password' parameter")
	// }

	token, err := utils.GenerateJWT(username)
	if err != nil {
		return nil, errors.New("Failed to login")
	}

	response := map[string]interface{}{"jwt_token": token}

	return response, nil
}

func IsUserExist(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoo.ClientUser.Database("userdb").Collection("users")
	filter := bson.D{{Key: "username", Value: username}}

	var result User

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, errors.New("Username already exist")
	}
}
