package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDatabaseConnection(t *testing.T) {
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/%s?retryWrites=true&w=majority", os.Getenv("WILHELMIINA_SERVER_USERNAME"), os.Getenv("WILHELMIINA_SERVER_PASSWORD"), os.Getenv("WILHELMIINA_SERVER_CLUSTER_NAME"), os.Getenv("WILHELMIINA_SERVER_DATABASE_NAME"))
	t.Log(connStr)
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		t.Error(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		t.Error(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		t.Error(err)
	}
	client.Disconnect(context.TODO())
}
