package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) ListEvents(ctx context.Context, title string, offset int64, limit int64) ([]types.Event, error) {
	filter := bson.M{}
	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	events := make([]types.Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	for i := range events {
		events[i].Address = events[i].Location.Address
	}

	return events, nil
}
