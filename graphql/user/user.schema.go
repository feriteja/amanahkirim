package user

import (
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

var ProfileField = &graphql.Field{
	Type:    ProfileType,
	Args:    ProfileArgs,
	Resolve: NewProfile,
}
