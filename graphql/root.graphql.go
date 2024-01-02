package graphql

import (
	"amanahkirim/graphql/product"
	"amanahkirim/graphql/user"

	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user":    user.UserField,
			"users":   user.UsersField,
			"login":   user.LoginField,
			"product": product.ProductField,
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": user.RegisterField,
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)
