package product

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func GetProduct(p graphql.ResolveParams) (interface{}, error) {
	user, _ := p.Context.Value("user").(string)

	fmt.Printf("token: %s", user)
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
