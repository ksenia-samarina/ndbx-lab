package cassandra

import (
	"context"
	"fmt"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Storage struct {
	session         *gocql.Session
	keyspace        string
	likeTableName   string
	reviewTableName string
}

func (s *Storage) InitSchema(ctx context.Context, keyspace string) error {
	if keyspace != "" {
		createKeyspaceQuery := fmt.Sprintf(`
			CREATE KEYSPACE IF NOT EXISTS %s 
			WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`,
			keyspace,
		)
		if err := s.session.Query(createKeyspaceQuery).ExecContext(ctx); err != nil {
			return fmt.Errorf("failed to create keyspace: %w", err)
		}
	}

	createLikeTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			event_id text,
			created_by text,
			like_value tinyint,
			created_at timestamp,
			PRIMARY KEY (event_id, created_by)
		);`,
		keyspace, s.likeTableName,
	)
	createReviewTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id uuid,
			event_id text,
			created_by text,
			rating tinyint,
			comment text,
			created_at timestamp,
			updated_at timestamp,
			PRIMARY KEY (event_id, id) 
		);`,
		keyspace, s.reviewTableName,
	)

	if err := s.session.Query(createLikeTableQuery).ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to create table %s: %w", s.likeTableName, err)
	}
	if err := s.session.Query(createReviewTableQuery).ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to create table %s: %w", s.reviewTableName, err)
	}

	createIndexQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS ON %s.%s (like_value);`,
		keyspace, s.likeTableName,
	)
	if err := s.session.Query(createIndexQuery).ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to create secondary index on like_value: %w", err)
	}

	createIndexQuery = fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS ON %s.%s (rating);`,
		keyspace, s.reviewTableName,
	)
	if err := s.session.Query(createIndexQuery).ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to create secondary index on rating: %w", err)
	}

	return nil
}

func NewCassandraStorage(session *gocql.Session, keyspace string, likeTableName string, reviewTableName string) *Storage {
	return &Storage{
		session:         session,
		keyspace:        keyspace,
		likeTableName:   likeTableName,
		reviewTableName: reviewTableName,
	}
}
