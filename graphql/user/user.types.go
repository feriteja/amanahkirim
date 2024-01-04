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

var AddSellerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AddSeller",
		Fields: graphql.Fields{
			"jwt_token": &graphql.Field{Type: graphql.String},
		},
	},
)

var AddSellerArgs = graphql.FieldConfigArgument{
	"name": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"phone_numbers": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"city": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"province": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"country": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"postal_code": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"profile_picture_url": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var AddBuyerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AddBuyer",
		Fields: graphql.Fields{
			"jwt_token": &graphql.Field{Type: graphql.String},
		},
	},
)

var AddBuyerArgs = graphql.FieldConfigArgument{
	"name": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"phone_numbers": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"city": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"province": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"regency": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"country": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"postal_code": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"date_of_birth": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"profile_picture_url": &graphql.ArgumentConfig{
		Type: graphql.String,
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
