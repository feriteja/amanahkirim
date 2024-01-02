package utils

import (
	"context"
	"net/http"
	"regexp"
)

var AuthenticatedQueries = map[string]struct{}{
	"product": {},
	"Query2":  {},
	// Add other query names as needed
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// Validate and extract user info from the token
		// Implement your token validation logic here
		// If validation fails, return an error and stop the request processing

		// For example:
		if token == "" || !IsValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If validation is successful, set the user context in the request context
		user := ExtractUserInfoFromToken(token)
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

func ExtractUserInfoFromToken(token string) string {
	// Implement your logic to extract user info from the token
	// Return the user info
	return "user123"
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
