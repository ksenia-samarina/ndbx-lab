package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetUserEventsByUserID(ctx context.Context, userID string) ([]model.Event, error) {
	return d.userStorage.GetUserEventsByUserID(ctx, userID)
}
