package reactions

import (
	"context"
	"fmt"
	"log"
	"samarina/ndbx/internal/model"
)

func (d *Domain) SetReaction(ctx context.Context, reaction model.Reaction) error {
	if reaction.LikeValue != 1 && reaction.LikeValue != -1 {
		return fmt.Errorf("invalid reaction value: %d", reaction.LikeValue)
	}

	event, err := d.eventsStorage.GetEventByID(ctx, reaction.EventID)
	if err != nil {
		return fmt.Errorf("find event: %w", err)
	}

	if err := d.reactionsStorage.SetReaction(ctx, reaction); err != nil {
		return fmt.Errorf("storage set reaction: %w", err)
	}

	field := "likes"
	if reaction.LikeValue == -1 {
		field = "dislikes"
	}

	err = d.reactionsCache.IncrementCounter(ctx, event.Title, field, 1, d.likeTTL)
	if err != nil {
		log.Printf("failed to increment redis counter: %v", err)
	}

	return nil
}
