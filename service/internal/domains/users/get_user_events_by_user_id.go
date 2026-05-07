package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetUserEventsByUserID(ctx context.Context, userID string, filter model.EventFilter) ([]model.Event, error) {
	return d.eventStorage.GetUserEventsByUserID(ctx, userID, filter)
}
