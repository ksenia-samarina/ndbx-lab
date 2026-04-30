package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetUserByUserID(ctx context.Context, id string) (model.User, error) {
	return d.userStorage.GetUserByUserID(ctx, id)
}
