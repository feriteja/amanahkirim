package user

import (
	"github.com/graphql-go/graphql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var Users = map[int]*User{
	1: {ID: 1, Name: "John Doe"},
	2: {ID: 2, Name: "Jane Bone"},
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
