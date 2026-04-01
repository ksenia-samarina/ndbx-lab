package mongo

import (
	"context"
	"errors"
	"samarina/ndbx/internal/domains/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	var user types.User

	filter := bson.M{"username": username}
	err := s.collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
