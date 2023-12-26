package main

import (
	_ "amanahkirim/db/mongoo"
	"encoding/json"
	"fmt"
	"net/http"

	qraphqlc "amanahkirim/graphql"
	"amanahkirim/graphql/field"

	"github.com/graphql-go/graphql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "John Doe"},
	{ID: 2, Name: "Jane Doe"},
}

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: qraphqlc.UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, user := range users {
							if user.ID == id {
								return user, nil
							}
						}
					}
					return nil, nil
				},
			},
			"users": &graphql.Field{
				Type: graphql.NewList(qraphqlc.UserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
		},
	},
)

// var schema, _ = graphql.NewSchema(
// 	graphql.SchemaConfig{
// 		Query: queryType,
// 	},
// )

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
		Schema:        field.UserSchema,
		RequestString: query,
	})

	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
