package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetUsers(ctx context.Context, id, name string, limit, offset uint64) ([]model.User, error) {
	return d.userStorage.GetUsers(ctx, id, name, limit, offset)
}
