package events

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	db               *mongo.Database
	eventsCollection *mongo.Collection
	usersCollection  *mongo.Collection
}

func (s *Storage) CreateIndexes(ctx context.Context, models []mongo.IndexModel) error {
	_, err := s.eventsCollection.Indexes().CreateMany(ctx, models)
	return err
}

func NewStorage(db *mongo.Database, eventsCollectionName string, usersCollectionName string) *Storage {
	return &Storage{
		db:               db,
		eventsCollection: db.Collection(eventsCollectionName),
		usersCollection:  db.Collection(usersCollectionName),
	}
}
