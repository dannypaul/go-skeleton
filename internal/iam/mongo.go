package iam

import (
	"github.com/dannypaul/go-skeleton/config"
	"github.com/dannypaul/go-skeleton/internal/driver/platform/mongo"
)

const UserCollectionName = "users"

type mongoUserRepo struct {
	mongo.Collection
}

func NewMongoUserRepo(client *mongo.Client) (UserRepo, error) {
	conf, err := config.Get()
	if err != nil {
		return nil, err
	}

	collection := mongo.Collection{Collection: client.Database(conf.MongoDbName).Collection(UserCollectionName)}

	return mongoUserRepo{collection}, err
}

const ChallengeCollectionName = "challenges"

type mongoChallengeRepo struct {
	mongo.Collection
}

func NewMongoChallengeRepo(client *mongo.Client) (ChallengeRepo, error) {
	conf, err := config.Get()
	if err != nil {
		return nil, err
	}

	collection := mongo.Collection{Collection: client.Database(conf.MongoDbName).Collection(ChallengeCollectionName)}

	return mongoChallengeRepo{collection}, err
}
