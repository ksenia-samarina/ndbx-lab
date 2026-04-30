package events

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) RegisterEvent(ctx context.Context, createdBy string, event model.Event) (string, error) {
	event.CreatedBy = createdBy
	event.CreatedAt = time.Now().Format(time.RFC3339)
	res, err := s.collection.InsertOne(ctx, &event)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
