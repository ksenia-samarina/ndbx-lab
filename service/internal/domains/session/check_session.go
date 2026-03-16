package session

import (
	"context"
)

func (d *Domain) CheckSessionExist(ctx context.Context, id Id) (bool, error) {
	return d.Storage.CheckExistenceSession(ctx, id)
}
