package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Data map[string]interface{} `json:"data"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func GenerateJWT(jwtData interface{}) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	expirationEnv := os.Getenv("JWT_EXPIRED")

	days, err := strconv.Atoi(expirationEnv)
	if err != nil {
		return "", fmt.Errorf("Error parsing expiration duration: %v", err)
	}
	expiration := time.Duration(days) * 24 * time.Hour

	dataMap := make(map[string]interface{})
	val := reflect.ValueOf(jwtData)
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			dataMap[field.Name] = val.Field(i).Interface()
		}
	}

	claims := Claims{
		Data: dataMap,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

	return tokenString, nil
}

func ComparePassword(hashedPassword []byte, inputPassword string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(inputPassword))
}
