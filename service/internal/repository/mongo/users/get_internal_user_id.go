package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) GetInternalUserID(ctx context.Context, username string) (primitive.ObjectID, error) {
	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	filter := bson.M{"username": username}

	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.ID, nil
}
