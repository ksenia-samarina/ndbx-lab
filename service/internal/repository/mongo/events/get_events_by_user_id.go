package events

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) GetEventsByUserID(ctx context.Context, createdBy string) ([]model.Event, error) {
	objID, _ := primitive.ObjectIDFromHex(createdBy)
	filter := bson.M{"created_by": objID}

	cursor, err := s.eventsCollection.Find(ctx, filter)
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
