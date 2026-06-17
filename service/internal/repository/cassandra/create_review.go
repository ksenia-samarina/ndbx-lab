package cassandra

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (s *Storage) CreateReview(ctx context.Context, r model.Review) (string, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s.%s (id, event_id, created_by, comment, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		IF NOT EXISTS`, s.keyspace, s.reviewTableName,
	)

	reviewUUID, err := gocql.ParseUUID(r.ID)
	if err != nil {
		return "", fmt.Errorf("invalid uuid: %w", err)
	}

	cleanCreatedBy := r.CreatedBy

	applied, err := s.session.Query(query,
		reviewUUID,
		r.EventID,
		cleanCreatedBy,
		r.Comment,
		r.Rating,
		r.CreatedAt,
		r.CreatedAt,
	).MapScanCASContext(ctx, make(map[string]interface{}))

	if err != nil {
		return "", fmt.Errorf("cassandra query error: %w", err)
	}
	if !applied {
		return "", ErrAlreadyExists
	}
	return r.ID, nil
}
