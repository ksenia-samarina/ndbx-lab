package users

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	filter := bson.M{"username": username}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
