package reactions

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) EnrichEventsWithReactions(ctx context.Context, events []model.Event, includeReactions bool) ([]model.Event, error) {
	if !includeReactions {
		return events, nil
	}

	for i := range events {
		defaultCounters := &model.ReactionCounters{Likes: 0, Dislikes: 0}

		counters, err := d.reactionsCache.GetCounters(ctx, events[i].Title)
		if err != nil {
			allEvents, dbErr := d.eventsStorage.GetEvents(ctx, model.EventFilter{
				Title: events[i].Title,
			})

			if dbErr != nil || len(allEvents) == 0 {
				allEvents = []model.Event{events[i]}
			}

			totalCounters := &model.ReactionCounters{Likes: 0, Dislikes: 0}

			for _, ev := range allEvents {
				c, cassErr := d.reactionsStorage.GetReactionCounters(ctx, ev.ID.Hex())
				if cassErr == nil && c != nil {
					totalCounters.Likes += c.Likes
					totalCounters.Dislikes += c.Dislikes
				}
			}

			counters = totalCounters

			_ = d.reactionsCache.SetCounters(ctx, events[i].Title, counters, d.likeTTL)
		}

		if counters == nil {
			counters = defaultCounters
		}

		events[i].Reactions = counters
	}

	return events, nil
}
