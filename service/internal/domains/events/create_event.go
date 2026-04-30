package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) CreateEvent(ctx context.Context, createdBy string, event *types.Event) (string, error) {
	id, err := d.eventsStorage.CreateEvent(ctx, createdBy, event)
	if err != nil {
		return "", err
	}
	return id, nil
}
