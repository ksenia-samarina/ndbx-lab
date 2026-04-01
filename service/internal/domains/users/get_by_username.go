package users

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	return d.userStorage.GetByUsername(ctx, username)
}
