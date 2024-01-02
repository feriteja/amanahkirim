package user

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.Int},
			"name": &graphql.Field{Type: graphql.String},
		},
	},
)

var CredentialType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Credential",
		Fields: graphql.Fields{
			"jwt_token": &graphql.Field{Type: graphql.String},
		},
	},
)

var RegisterType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Register",
		Fields: graphql.Fields{
			"jwt_token": &graphql.Field{Type: graphql.String},
		},
	},
)

var RegisterArgs = graphql.FieldConfigArgument{
	"username": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"password": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"confirm_password": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
}

var UserArgs = graphql.FieldConfigArgument{
	"name": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var LoginArgs = graphql.FieldConfigArgument{
	"username": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"password": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}
