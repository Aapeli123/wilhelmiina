package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const password = os.Getenv("WILHELMIINA_SERVER_PASSWORD")
const connStr = fmt.Sprintf("mongodb+srv://admin:%s@wilhelmiinatest.8tutg.mongodb.net/wilhelmiinatest?retryWrites=true&w=majority", password)

// DbClient should be used for database operations
var DbClient *mongo.Client

// Init sets DbClient to the mongo client and tests connection.
// Panics on connection error, as the app would unusable
func Init() {
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		panic(err)
	}

	DbClient = client // Assing the global DbClient
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = DbClient.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = DbClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection succesful")
}

// Close closes the connection to database
func Close() {
	DbClient.Disconnect(context.TODO())
	fmt.Println("Database connection closed")
}
