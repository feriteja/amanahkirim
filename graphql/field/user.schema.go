package field

import (
	"github.com/graphql-go/graphql"

	qraphqlc "amanahkirim/graphql"
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

var UserSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)
