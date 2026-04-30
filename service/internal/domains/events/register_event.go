package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) RegisterEvent(ctx context.Context, createdBy string, event model.Event) (string, error) {
	return d.eventsStorage.RegisterEvent(ctx, createdBy, event)
}
