package product

import (
	"github.com/graphql-go/graphql"
)

type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Seller int    `json:"seller"`
}

var Products = map[int]*Product{
	1: {ID: 1, Name: "Ultramilk", Seller: 1},
	2: {ID: 2, Name: "FlagBear", Seller: 1},
	3: {ID: 3, Name: "Meatball", Seller: 2},
	4: {ID: 4, Name: "TunaFish", Seller: 2},
	5: {ID: 5, Name: "Chicken", Seller: 2},
}

var ProductField = &graphql.Field{
	Type:    ProductType,
	Args:    ProductArgs,
	Resolve: GetProduct,
}
