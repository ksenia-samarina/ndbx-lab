package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) CreateEvent(ctx context.Context, createdBy string, event *types.Event) (string, error) {
	event.CreatedBy = createdBy
	event.CreatedAt = time.Now().Format(time.RFC3339)
	res, err := s.collection.InsertOne(ctx, &event)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
