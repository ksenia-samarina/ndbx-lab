package users

import (
	"context"
	"log"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetUserEventsByUserID(ctx context.Context, userID string) ([]model.Event, error) {
	filter := bson.M{"created_by": userID}

	cursor, err := s.collection.Find(ctx, filter)
	log.Printf("curssor filters: %v", filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	events := make([]model.Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	log.Printf("curssor events: %v", events)
	return events, nil
}
