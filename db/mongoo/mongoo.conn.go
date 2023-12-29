package mongoo

import (
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ClientUser    *mongo.Client
	ClientProduct *mongo.Client
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	initUserDb()
	initProductDb()
}
