package cassandra

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (s *Storage) UpdateReview(ctx context.Context, r model.Review) error {
	query := fmt.Sprintf(`
		UPDATE %s.event_reviews 
		SET comment = ?, rating = ?, updated_at = ?
		WHERE event_id = ? AND id = ?`, s.keyspace,
	)

	reviewUUID, err := gocql.ParseUUID(r.ID)
	if err != nil {
		return fmt.Errorf("invalid review uuid for update: %w", err)
	}

	err = s.session.Query(query,
		r.Comment,
		r.Rating,
		time.Now().UTC(),
		r.EventID,
		reviewUUID,
	).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("cassandra update review query error: %w", err)
	}

	return nil
}
