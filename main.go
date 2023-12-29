package main

import (
	_ "amanahkirim/db/mongoo"
	"encoding/json"
	"fmt"
	"net/http"

	graphqlc "amanahkirim/graphql"

	"github.com/graphql-go/graphql"
)

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var query string

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var request map[string]interface{}
		err := decoder.Decode(&request)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		query, _ = request["query"].(string)
	} else {
		query = r.URL.Query().Get("query")
	}

	result := graphql.Do(graphql.Params{
		Schema:        graphqlc.Schema,
		RequestString: query,
	})

	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
