package product

import (
	"amanahkirim/graphql/user"

	"github.com/graphql-go/graphql"
)

var ProductType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.Int},
			"name": &graphql.Field{Type: graphql.String},
			"sellerInfo": &graphql.Field{Type: user.UserType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// if product, ok := p.Source.(*Product); ok {
					// 	log.Print(product.Seller)
					// 	return mongoo.User[product.Seller], nil
					// }
					return nil, nil
				},
			},
		},
	},
)

var ProductArgs = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}
