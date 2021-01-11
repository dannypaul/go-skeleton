package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/dannypaul/go-skeleton/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/rs/zerolog/log"
)

type Client struct {
	*mongo.Client
}

// Connect ...
func Connect(ctx context.Context) *Client {
	conf, _ := config.Get()
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.MongoURI))
	if err != nil {
		log.Fatal().Err(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("db connection error: %s", err))
	}

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(fmt.Errorf("db ping error: %s", err))
	}

	log.Info().Msg("Successfully Connected to MongoDB")

	return &Client{client}
}

// Disconnect ...
func Disconnect(client *Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Info().Err(fmt.Errorf("disconnection from MongoDB not successful %w", err))
	}
	log.Info().Msg("Disconnected from MongoDB successfully")
}
