package reactions

import (
	"context"
	"log"
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
			dbCounters, dbErr := d.reactionsStorage.GetReactionCounters(ctx, events[i].ID.Hex())
			if dbErr != nil {
				log.Printf("Warning: failed to get reactions from cassandra for event %s: %v", events[i].ID.Hex(), dbErr)
				events[i].Reactions = defaultCounters
				continue
			}

			counters = dbCounters
			_ = d.reactionsCache.SetCounters(ctx, events[i].Title, counters, d.likeTTL)
		}

		if counters == nil {
			counters = defaultCounters
		}

		events[i].Reactions = counters
	}

	return events, nil
}
