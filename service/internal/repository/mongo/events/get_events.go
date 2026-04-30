package events

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetEvents(ctx context.Context, filters model.EventFilter) ([]model.Event, error) {
	filter := bson.M{}
	if filters.Title != "" {
		filter["title"] = bson.M{"$regex": filters.Title, "$options": "i"}
	}

	findOptions := options.Find()
	if filters.Limit > 0 {
		findOptions.SetLimit(filters.Limit)
	}
	findOptions.SetSkip(filters.Offset)
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

	events := make([]model.Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}
