package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	var user types.User

	filter := bson.M{"username": username}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
