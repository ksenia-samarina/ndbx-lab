package login

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetByUsername(ctx context.Context, username string) (model.User, error) {
	return d.loginStorage.GetUserByUsername(ctx, username)
}
