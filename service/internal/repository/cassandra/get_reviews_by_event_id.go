package cassandra

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (s *Storage) GetReviewsByEventID(ctx context.Context, eventID string) ([]model.Review, error) {
	query := fmt.Sprintf(`
		SELECT id, event_id, created_by, comment, rating, created_at 
		FROM %s.%s 
		WHERE event_id = ?`, s.keyspace, s.reviewTableName,
	)

	scanner := s.session.Query(query, eventID).WithContext(ctx).Iter().Scanner()
	var reviews []model.Review

	for scanner.Next() {
		var r model.Review
		var reviewUUID gocql.UUID

		err := scanner.Scan(&reviewUUID, &r.EventID, &r.CreatedBy, &r.Comment, &r.Rating, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		r.ID = reviewUUID.String()
		reviews = append(reviews, r)
	}
	return reviews, scanner.Err()
}
