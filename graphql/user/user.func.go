package user

import (
	"fmt"

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

	user.Name = newName
	return user, nil
}

func GetAllUsers(p graphql.ResolveParams) (interface{}, error) {
	return Users, nil
}
