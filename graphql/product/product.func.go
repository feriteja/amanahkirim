package product

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

func GetProduct(p graphql.ResolveParams) (interface{}, error) {
	// user, _ := p.Context.Value("user").(string)

	userValue := p.Context.Value("user")
	userClaims, ok := userValue.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to assert the type of 'user'")
		return nil, fmt.Errorf("Failed to assert the type of 'user'")
	}

	// Access the "username" field
	username, exists := userClaims["username"].(string)
	if !exists {
		fmt.Println("Failed to get the 'username' field")
		return nil, fmt.Errorf("Failed to get the 'username' field")
	}

	fmt.Print(username)

	userID, ok := p.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' parameter")
	}

	product := Products[userID]
	if product != nil {
		return product, nil
	}

	return nil, fmt.Errorf("user with ID %d not found", userID)
}
