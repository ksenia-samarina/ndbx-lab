package users

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetUserEventsByUserID(ctx context.Context, userID string) ([]model.Event, error) {
	filter := bson.M{"created_by": userID}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	events := make([]model.Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}
