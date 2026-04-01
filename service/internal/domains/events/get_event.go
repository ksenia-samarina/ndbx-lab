package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) GetEvent(ctx context.Context, createdBy string) (*types.Event, error) {
	return d.eventsStorage.GetEvent(ctx, createdBy)
}
