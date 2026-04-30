package events

import (
	"context"
)

func (d *Domain) UpdateEventsLocationCity(ctx context.Context, id string, city string) error {
	return d.eventsStorage.UpdateEventsLocationCity(ctx, id, city)
}
