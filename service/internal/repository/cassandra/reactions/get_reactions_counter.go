package reactions

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
)

func (s *Storage) GetReactionCounters(ctx context.Context, eventID string) (*model.ReactionCounters, error) {
	query := `SELECT like_value FROM ` + s.keyspace + `.` + s.tableName + ` WHERE event_id = ?`

	scanner := s.session.Query(query, eventID).IterContext(ctx).Scanner()

	counters := &model.ReactionCounters{
		Likes:    0,
		Dislikes: 0,
	}

	var likeValue int8

	for scanner.Next() {
		if err := scanner.Scan(&likeValue); err != nil {
			return nil, fmt.Errorf("scan reaction row: %w", err)
		}
		switch likeValue {
		case 1:
			counters.Likes++
		case -1:
			counters.Dislikes++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cassandra iterator error: %w", err)
	}

	return counters, nil
}
