package events

import (
	"context"
	"errors"
	"log"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetEventByID(ctx context.Context, id string) (model.Event, error) {
	objID, _ := primitive.ObjectIDFromHex(id)

	var event model.Event
	err := s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	log.Printf("got event in GetEventById: %v, %s", event, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("got error event in GetEventById: %v, %s", event, id)
		return model.Event{}, ErrEventNotExist
	}
	if err != nil {
		return model.Event{}, err
	}
	return event, nil
}
