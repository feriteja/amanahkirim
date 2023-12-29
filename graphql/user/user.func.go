package user

import (
	"amanahkirim/db/mongoo"
	"amanahkirim/graphql/user/utils"
	rootUtils "amanahkirim/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

func GetUser(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := p.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' parameter")
	}

	user := Users[userID]
	if user != nil {
		return user, nil
	}

	return nil, fmt.Errorf("user with ID %d not found", userID)
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

func CreateUser(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	username := p.Args["username"].(string)
	password := p.Args["password"].(string)

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, &rootUtils.AppResponse{Code: http.StatusBadRequest, Message: "Failed to encrypt password"}
	}

	user := mongoo.User{
		Username: username,
		Password: hashPassword,
	}

	collection := mongoo.ClientUser.Database("userdb").Collection("users")

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, &rootUtils.AppResponse{Code: http.StatusBadRequest, Message: "Failed to create new user"}
	}

	response := map[string]interface{}{"username": username}

	return response, nil

}
