package cassandra

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (s *Storage) SetReaction(ctx context.Context, reaction model.Reaction) error {
	query := `
		INSERT INTO ` + s.keyspace + `.` + s.likeTableName + ` (event_id, created_by, like_value, created_at)
		VALUES (?, ?, ?, ?)
	`

	err := s.session.Query(
		query,
		reaction.EventID,
		reaction.CreatedBy,
		reaction.LikeValue,
		reaction.CreatedAt,
	).ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}
