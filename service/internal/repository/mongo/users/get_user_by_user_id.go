package users

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) GetUserByUserID(ctx context.Context, id string) (model.User, error) {
	objID, _ := primitive.ObjectIDFromHex(id)

	var user model.User
	err := s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
