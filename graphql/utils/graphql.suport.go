package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

type JwtData struct {
	ID       string  `json:"_id"`
	BuyerID  *string `json:"buyer_id"`
	SellerID *string `json:"seller_id"`
}

func GetUserFromContext(p graphql.ResolveParams) (*JwtData, error) {
	userValue := p.Context.Value("user")
	fmt.Println(userValue)

	// Check if userValue is not nil
	if userValue == nil {
		log.Println("'user' value is nil")
		return nil, errors.New("'user' value is nil")
	}

	userClaims, ok := userValue.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to assert the type of 'user'")
		return nil, errors.New("Failed to assert the type of 'user")
	}

	// Check if "data" key is present in userClaims
	dataClaim, ok := userClaims["data"]
	if !ok {
		log.Println("'data' key is not present in 'user' claims")
		return nil, errors.New("'data' key is not present in 'user' claims")
	}

	// Type assertion for dataClaim
	data, ok := dataClaim.(map[string]interface{})
	if !ok {
		log.Println("Failed to assert the type of 'data'")
		return nil, errors.New("Failed to assert the type of 'data'")
	}

	// Create JwtData from the map
	jwtData := JwtData{
		ID:       getString(data, "ID"),
		BuyerID:  getStringPtr(data, "BuyerID"),
		SellerID: getStringPtr(data, "SellerID"),
	}

	return &jwtData, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getStringPtr(data map[string]interface{}, key string) *string {
	if val, ok := data[key].(string); ok {
		return &val
	}
	return nil
}
