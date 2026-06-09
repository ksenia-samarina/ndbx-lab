package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEventsByUserID(ctx context.Context, createdBy string) ([]model.Event, error) {
	return d.eventsStorage.GetEventsByUserID(ctx, createdBy)
}
