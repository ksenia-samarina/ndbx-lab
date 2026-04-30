package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) ListEvents(ctx context.Context, title string, offset int64, limit int64) ([]types.Event, error) {
	return d.eventsStorage.ListEvents(ctx, title, offset, limit)
}
