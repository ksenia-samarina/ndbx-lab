package events

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) UpdateEventsLocationCity(ctx context.Context, id string, city string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": bson.M{"location.city": city},
	}
	_, err := s.eventsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
