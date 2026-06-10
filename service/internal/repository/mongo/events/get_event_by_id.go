package events

import (
	"context"
	"errors"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetEventByID(ctx context.Context, id string) (model.Event, error) {
	objID, _ := primitive.ObjectIDFromHex(id)

	var event model.Event
	err := s.eventsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Event{}, ErrEventNotExist
	}
	if err != nil {
		return model.Event{}, err
	}
	return event, nil
}
