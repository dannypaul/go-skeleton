package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Copier struct {
	result *mongo.SingleResult
}

func (c Copier) Copy(destination interface{}) error {
	if c.result.Err() != nil {
		return c.result.Err()
	}
	return c.result.Decode(destination)
}

type ListCopier struct {
	cursor *mongo.Cursor
}

func (l ListCopier) CopyAll(ctx context.Context, destination interface{}) error {
	if l.cursor.Err() != nil {
		return l.cursor.Err()
	}
	return l.cursor.All(ctx, destination)
}
