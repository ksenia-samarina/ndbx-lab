package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) CreateEvent(ctx context.Context, createdBy string, event *types.Event) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(createdBy)
	if err != nil {
		return "", err
	}
	event.CreatedBy = objectID
	event.Location.Address = event.Address
	event.CreatedAt = time.Now()
	res, err := s.collection.InsertOne(ctx, &event)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
