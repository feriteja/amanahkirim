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
			http.Error(w, "Error decoding request bodyss", http.StatusBadRequest)
			return
		}
		query, _ = request["query"].(string)
	} else {
		query = r.URL.Query().Get("query")
	}

	// user, _ := r.Context().Value("user").(string)
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

	json.NewEncoder(w).Encode(result)

}

func main() {
	handler := utils.AuthenticationMiddleware(http.HandlerFunc(graphqlHandler))
	http.Handle("/graphql", handler)
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
