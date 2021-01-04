package mongo

import (
	"context"
	"errors"

	"github.com/dannypaul/go-skeleton/internal/exception"
	"github.com/dannypaul/go-skeleton/internal/primitive"
	"github.com/dannypaul/go-skeleton/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func (c Collection) Set(ctx context.Context, filters []repository.Filter, Key string, value interface{}) error {
	match := bson.D{}
	for _, f := range filters {
		match = append(match, bson.E{Key: f.Key, Value: f.Value})
	}
	update := bson.D{{"$set", bson.D{{Key, value}}}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) SetAll(ctx context.Context, filters []repository.Filter, keyValues []repository.KeyValue) error {
	match := bson.D{}
	for _, f := range filters {
		match = append(match, bson.E{Key: f.Key, Value: f.Value})
	}

	setters := bson.D{}
	for _, kv := range keyValues {
		setters = append(setters, bson.E{Key: kv.Key, Value: kv.Value})
	}
	update := bson.D{{"$set", setters}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) SetById(ctx context.Context, id primitive.Id, Key string, Value interface{}) error {
	match := bson.D{{"_id", id}}

	setter := bson.D{{Key, Value}}
	update := bson.D{{"$set", setter}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) SetAllById(ctx context.Context, id primitive.Id, keyValues []repository.KeyValue) error {
	match := bson.D{{"_id", id}}

	setters := bson.D{}
	for _, kv := range keyValues {
		setters = append(setters, bson.E{Key: kv.Key, Value: kv.Value})
	}
	update := bson.D{{"$set", setters}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) UnSet(ctx context.Context, filters []repository.Filter, Key string) error {
	match := bson.D{}
	for _, f := range filters {
		match = append(match, bson.E{Key: f.Key, Value: f.Value})
	}
	update := bson.D{{"$unset", bson.D{{Key, ""}}}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) Patch(ctx context.Context, id primitive.Id, patches []repository.Patch) error {
	match := bson.D{{"_id", id}}

	patch := bson.D{}
	for _, p := range patches {
		patch = append(patch, bson.E{Key: p.Action, Value: bson.D{{p.Key, p.Value}}})
	}

	res, err := c.UpdateOne(ctx, match, patch)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) IncrementById(ctx context.Context, id primitive.Id, Key string, incrementBy int) error {
	match := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{Key, incrementBy}}}}

	res, err := c.UpdateOne(ctx, match, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrNotFound
	}
	return err
}

func (c Collection) Count(ctx context.Context, filters []repository.Filter) (int64, error) {
	match := bson.D{}
	for _, f := range filters {
		match = append(match, bson.E{Key: f.Key, Value: f.Value})
	}
	return c.CountDocuments(ctx, match)
}

func (c Collection) Create(ctx context.Context, model interface{}) (repository.Copier, error) {
	doc, err := bson.Marshal(model)
	if err != nil {
		return nil, err
	}

	result, err := c.InsertOne(ctx, doc)
	if err != nil {
		var e mongo.WriteException
		if errors.As(err, &e) {
			for _, we := range e.WriteErrors {
				if we.Code == 11000 {
					return nil, exception.ErrConflict
				}
			}
		}
		return nil, err
	}

	singleResult := c.FindOne(ctx, bson.M{"_id": result.InsertedID})
	return Copier{singleResult}, singleResult.Err()
}

func (c Collection) FindById(ctx context.Context, id primitive.Id) (repository.Copier, error) {
	singleResult := c.FindOne(ctx, bson.M{"_id": id})
	if singleResult.Err() != nil && errors.Is(singleResult.Err(), mongo.ErrNoDocuments) {
		return nil, exception.ErrNotFound
	}
	return Copier{singleResult}, singleResult.Err()
}

func (c Collection) Delete(ctx context.Context, id primitive.Id) (int64, error) {
	deleteResult, err := c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func (c Collection) Add(ctx context.Context, filters []repository.Filter, key string) (float64, error) {
	matchConditions := bson.D{}
	for _, f := range filters {
		matchConditions = append(matchConditions, bson.E{Key: f.Key, Value: f.Value})
	}
	match := bson.D{{"$match", matchConditions}}
	group := bson.D{{"$group", bson.D{{"total", bson.D{{"$sum", "$" + key}}}}}}

	cursor, err := c.Aggregate(ctx, mongo.Pipeline{match, group})
	if err != nil {
		return 0, err
	}

	var res []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &res); err != nil {
		return 0, err
	}
	return res[0].Total, err

}

func (c Collection) FindSingle(ctx context.Context, filters []repository.Filter) (repository.Copier, error) {
	match := bson.D{}
	for _, f := range filters {
		match = append(match, bson.E{Key: f.Key, Value: f.Value})
	}

	singleResult := c.FindOne(ctx, match)
	if singleResult.Err() != nil && errors.Is(singleResult.Err(), mongo.ErrNoDocuments) {
		return nil, exception.ErrNotFound
	}
	return Copier{singleResult}, singleResult.Err()
}

func (c Collection) Replace(ctx context.Context, id primitive.Id, value interface{}) (repository.Copier, error) {
	doc, err := bson.Marshal(value)
	if err != nil {
		return nil, err
	}

	after := options.After
	replaceOptions := options.FindOneAndReplaceOptions{ReturnDocument: &after}
	singleResult := c.FindOneAndReplace(ctx, bson.M{"_id": id}, doc, &replaceOptions)
	return Copier{singleResult}, singleResult.Err()
}
