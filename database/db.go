package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbClient should be used for database operations
var DbClient *mongo.Client

// Init sets DbClient to the mongo client and tests connection.
// Panics on connection error, as the app would unusable
func Init() {
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/%s?retryWrites=true&w=majority", os.Getenv("WILHELMIINA_SERVER_USERNAME"), os.Getenv("WILHELMIINA_SERVER_PASSWORD"), os.Getenv("WILHELMIINA_SERVER_CLUSTER_NAME"), os.Getenv("WILHELMIINA_SERVER_DATABASE_NAME"))
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
}

// Close closes the connection to database
func Close() {
	DbClient.Disconnect(context.TODO())
}
