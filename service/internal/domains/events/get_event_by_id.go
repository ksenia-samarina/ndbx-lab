package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEventByID(ctx context.Context, id string) (model.Event, error) {
	return d.eventsStorage.GetEventByID(ctx, id)
}
