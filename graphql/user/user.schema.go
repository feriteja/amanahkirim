package user

import (
	"amanahkirim/graphql/user/seller"

	"github.com/graphql-go/graphql"
)

// type User struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty"`
// 	Username string             `bson:"username" unique:"true"`
// 	Password string             `bson:"password"`
// }

var LoginField = &graphql.Field{
	Type:    CredentialType,
	Args:    LoginArgs,
	Resolve: Login,
}
var RegisterField = &graphql.Field{
	Type:    RegisterType,
	Args:    RegisterArgs,
	Resolve: RegisterUser,
}

var AddBuyerField = &graphql.Field{
	Type:    AddBuyerType,
	Args:    AddBuyerArgs,
	Resolve: NewBuyer,
}

var AddSellerField = &graphql.Field{
	Type:    AddSellerType,
	Args:    AddSellerArgs,
	Resolve: seller.NewSeller,
}
