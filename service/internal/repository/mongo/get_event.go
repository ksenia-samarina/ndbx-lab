package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) GetEvent(ctx context.Context, createdBy string) (*types.Event, error) {
	objectID, err := primitive.ObjectIDFromHex(createdBy)
	if err != nil {
		return nil, err
	}

	var event types.Event
	filter := bson.M{"createdBy": objectID}

	err = s.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
