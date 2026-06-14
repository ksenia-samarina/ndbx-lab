package reactions

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
)

func (d *Domain) SetReaction(ctx context.Context, reaction model.Reaction) error {
	if reaction.LikeValue != 1 && reaction.LikeValue != -1 {
		return fmt.Errorf("invalid reaction value: %d", reaction.LikeValue)
	}

	currentEvent, err := d.eventsStorage.GetEventByID(ctx, reaction.EventID)
	if err != nil {
		return fmt.Errorf("find event error: %w", err)
	}

	if err := d.reactionsStorage.SetReaction(ctx, reaction); err != nil {
		return fmt.Errorf("storage set reaction error: %w", err)
	}

	allEvents, err := d.eventsStorage.GetEvents(ctx, model.EventFilter{
		Title: currentEvent.Title,
	})
	if err != nil || len(allEvents) == 0 {
		allEvents = []model.Event{currentEvent}
	}

	totalCounters := &model.ReactionCounters{Likes: 0, Dislikes: 0}
	for _, ev := range allEvents {
		idSearch := ev.ID.Hex()

		counters, dbErr := d.reactionsStorage.GetReactionCounters(ctx, idSearch)
		if dbErr == nil && counters != nil {
			totalCounters.Likes += counters.Likes
			totalCounters.Dislikes += counters.Dislikes
		}
	}

	err = d.reactionsCache.SetCounters(ctx, currentEvent.Title, totalCounters, d.likeTTL)
	if err != nil {
		return fmt.Errorf("failed to update redis cache: %w", err)
	}

	return nil
}
