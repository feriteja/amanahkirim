package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
)

var AuthenticatedQueries = map[string]struct{}{
	"add_buyer":  {},
	"add_seller": {},
	"add_user":   {},
	// Add other query names as needed
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var query string
		buf, err := io.ReadAll(r.Body) // handle the error
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = rdr2

		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(rdr1)
			r.Body.Close()
			var request map[string]interface{}
			if err := decoder.Decode(&request); err != nil {
				http.Error(w, "Error decoding request bodyss", http.StatusBadRequest)
				return
			}
			query, _ = request["query"].(string)
		} else {
			query = r.URL.Query().Get("query")
		}

		if query != "" {
			queryName := ExtractQueryName(query)
			if _, requiresAuth := AuthenticatedQueries[queryName]; !requiresAuth {
				next.ServeHTTP(w, r)
				return
			}
		}

		token := r.Header.Get("Authorization")

		// For example:
		if token == "" || !IsValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If validation is successful, set the user context in the request context
		user, err := ExtractUserInfoFromToken(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func IsValidToken(token string) bool {
	// Implement your token validation logic here
	// Return true if the token is valid, false otherwise
	return true
}

func ExtractUserInfoFromTokena(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Use the same secret key as used for signing
		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func ExtractUserInfoFromToken(tokenString string) (jwt.Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for validation
		return []byte(jwtSecret), nil
	})

	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}

	// Check if the token is valid

	return token.Claims, nil
}

func ExtractQueryName(query string) string {
	// Implement your logic to extract the query name from the GraphQL request
	// For simplicity, assume the query name is the first word after "query"

	re := regexp.MustCompile(`\b(mutation|query) (\w+) {[\s\n]*([^\s({]+)`)

	matches := re.FindStringSubmatch(query)

	if len(matches) > 0 {
		return matches[3]
	}
	return ""
}
