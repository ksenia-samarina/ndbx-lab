package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error) {
	events, err := d.eventsStorage.GetEvents(ctx, filter)
	return events, err
}
