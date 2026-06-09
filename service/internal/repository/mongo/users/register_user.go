package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (s *Storage) RegisterUser(ctx context.Context, user model.User) error {
	_, err := s.collection.InsertOne(ctx, user)
	return err
}
