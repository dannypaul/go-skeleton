package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dannypaul/go-skeleton/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	*mongo.Client
}

// Connect ...
func Connect(ctx context.Context) *Client {
	conf, _ := config.Get()
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("db connection error: %s", err))
	}

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(fmt.Errorf("db ping error: %s", err))
	}

	log.Println("Successfully Connected to MongoDB")

	return &Client{client}
}

// Disconnect ...
func Disconnect(client *Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Println("Disconnection from MongoDB not successful")
		log.Println(err)
	}
	log.Println("Disconnected from MongoDB successfully")
}
