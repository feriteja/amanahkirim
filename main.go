package main

import (
	_ "amanahkirim/db/mongoo"
	"amanahkirim/utils"
	"encoding/json"
	"fmt"
	"net/http"

	graphqlc "amanahkirim/graphql"

	"github.com/graphql-go/graphql"
)

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var query string

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var request map[string]interface{}
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		query, _ = request["query"].(string)
	} else {
		query = r.URL.Query().Get("query")
	}

	if query != "" {
		queryName := utils.ExtractQueryName(query)

		if _, requiresAuth := utils.AuthenticatedQueries[queryName]; requiresAuth {
			if r.Context().Value("user") == nil {

				http.Error(w, "Unauthorized: Missing Authorization token", http.StatusUnauthorized)
				return
			}
		}
	}

	result := graphql.Do(graphql.Params{
		Schema:        graphqlc.Schema,
		RequestString: query,
		Context:       r.Context(),
	})

	if len(result.Errors) > 0 {
		errList := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			errList[i] = err.Message
		}
		http.Error(w, fmt.Sprintf("GraphQL error(s): %v", errList), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	handler := utils.AuthenticationMiddleware(http.HandlerFunc(graphqlHandler))
	http.Handle("/graphql", handler)
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
