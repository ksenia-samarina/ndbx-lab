package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error) {
	return d.eventStorage.GetEvents(ctx, filter)
}
