package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error) {
	events, err := d.eventsStorage.GetEvents(ctx, filter)
	seen := make(map[string]bool)
	uniqueEvents := make([]model.Event, 0)

	for _, event := range events {
		if !seen[event.Title] {
			seen[event.Title] = true
			uniqueEvents = append(uniqueEvents, event)
		}
	}
	return uniqueEvents, err
}
