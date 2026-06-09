package users

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error) {
	events, err := d.eventStorage.GetEvents(ctx, filter)
	seen := make(map[string]struct{})
	uniqueEvents := make([]model.Event, 0)

	for _, event := range events {
		if _, exists := seen[event.Title]; !exists {
			seen[event.Title] = struct{}{}
			uniqueEvents = append(uniqueEvents, event)
		}
	}
	return uniqueEvents, err
}
