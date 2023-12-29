package user

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}

var Users = map[int]*User{
	1: {ID: primitive.NewObjectID(), Username: "John Doe", Password: ""},
}

var UserField = &graphql.Field{
	Type:    UserType,
	Args:    UserArgs,
	Resolve: GetUser,
}

var UsersField = &graphql.Field{
	Type:    graphql.NewList(UserType),
	Resolve: GetAllUsers,
}
