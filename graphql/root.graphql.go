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
			"login":   user.LoginField,
			"product": product.ProductField,
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register":   user.RegisterField,
			"add_buyer":  user.AddBuyerField,
			"add_seller": user.AddSellerField,
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)
