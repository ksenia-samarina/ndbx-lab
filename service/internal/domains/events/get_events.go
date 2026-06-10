package events

import (
	"context"
	"log"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error) {
	events, err := d.eventsStorage.GetEvents(ctx, filter)
	log.Printf("EVENTS: %+v", events)
	return events, err
}
