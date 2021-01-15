package mongo

import (
	"context"
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
		log.Fatal().Err(err).Msg("Error occurred while creating the MongoDB client")
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Error occurred while trying to connect to MongoDB")
	}

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("Error occurred while trying to ping MongoDB")
	}

	log.Info().Msg("Successfully Connected to MongoDB")

	return &Client{client}
}

// Disconnect ...
func Disconnect(client *Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error occurred when trying to disconnect from MongoDB")
	}
	log.Info().Msg("Disconnected from MongoDB successfully")
}
