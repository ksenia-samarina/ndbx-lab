package mongo

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (s *Storage) RegisterUser(ctx context.Context, user *types.User) error {
	_, err := s.collection.InsertOne(ctx, user)
	return err
}
