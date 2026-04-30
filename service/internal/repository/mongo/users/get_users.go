package users

import (
	"context"
	"samarina/ndbx/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetUsers(ctx context.Context, id, name string, limit, offset uint64) ([]model.User, error) {
	filter := bson.M{}

	if id != "" {
		objID, _ := primitive.ObjectIDFromHex(id)
		filter["_id"] = objID
	}
	if name != "" {
		filter["full_name"] = bson.M{"$regex": name, "$options": "i"}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]model.User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
