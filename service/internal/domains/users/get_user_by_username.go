package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	return d.userStorage.GetUserByUsername(ctx, username)
}
