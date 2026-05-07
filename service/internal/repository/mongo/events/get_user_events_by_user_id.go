package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (s *Storage) GetUserEventsByUserID(ctx context.Context, userID string, filter model.EventFilter) ([]model.Event, error) {
	//cursor, err := s.collection.Find(ctx, filter)
	//log.Printf("curssor filters: %v", filter)
	//if err != nil {
	//	return nil, err
	//}
	//defer cursor.Close(ctx)
	//
	//events := make([]model.Event, 0)
	//if err := cursor.All(ctx, &events); err != nil {
	//	return nil, err
	//}
	//log.Printf("curssor events: %v", events)
	events, _ := s.GetEvents(ctx, filter)
	seen := make(map[string]bool)
	uniqueEvents := make([]model.Event, 0)

	for _, event := range events {
		if !seen[event.Title] {
			seen[event.Title] = true
			uniqueEvents = append(uniqueEvents, event)
		}
	}
	return events, nil
}
