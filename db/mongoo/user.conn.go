package mongoo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initUserDb() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database credentials from environment variables
	username := os.Getenv("DB_USER_USERNAME")
	password := os.Getenv("DB_USER_PASSWORD")
	host := os.Getenv("DB_USER_HOST")
	port := os.Getenv("DB_USER_PORT")
	dbName := os.Getenv("DB_USER_NAME")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	clientOptions := options.Client().ApplyURI(uri)
	ClientUser, _ = mongo.Connect(context.Background(), clientOptions)

	// Check the connection
	err := ClientUser.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB user connected successfully!")
}
