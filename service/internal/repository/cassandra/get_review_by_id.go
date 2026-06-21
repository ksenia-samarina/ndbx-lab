package cassandra

import (
	"context"
	"errors"
	"fmt"
	"samarina/ndbx/internal/model"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (s *Storage) GetReviewByID(ctx context.Context, eventID, reviewID string) (*model.Review, error) {
	query := fmt.Sprintf(`
		SELECT id, event_id, created_by, comment, rating, created_at, updated_at 
		FROM %s.%s 
		WHERE event_id = ?`, s.keyspace, s.reviewTableName,
	)

	scanner := s.session.Query(query, eventID).IterContext(ctx).Scanner()

	targetUUID, err := gocql.ParseUUID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid review uuid format: %w", err)
	}

	for scanner.Next() {
		var r model.Review
		var reviewUUID gocql.UUID

		err := scanner.Scan(
			&reviewUUID,
			&r.EventID,
			&r.CreatedBy,
			&r.Comment,
			&r.Rating,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan review fields error: %w", err)
		}

		if reviewUUID == targetUUID {
			r.ID = reviewUUID.String()
			return &r, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("review not found")
}
