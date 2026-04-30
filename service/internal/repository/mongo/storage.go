package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func (s *Storage) CreateIndexes(ctx context.Context, models []mongo.IndexModel) error {
	_, err := s.collection.Indexes().CreateMany(ctx, models)
	return err
}

func NewStorage(db *mongo.Database, collectionName string) *Storage {
	return &Storage{
		db:         db,
		collection: db.Collection(collectionName),
	}
}
