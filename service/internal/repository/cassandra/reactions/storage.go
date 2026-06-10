package reactions

import (
	"context"
	"fmt"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Storage struct {
	session   *gocql.Session
	keyspace  string
	tableName string
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

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			event_id text,
			created_by text,
			like_value tinyint,
			created_at timestamp,
			PRIMARY KEY (event_id, created_by)
		);`,
		keyspace, s.tableName,
	)

	if err := s.session.Query(createTableQuery).ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to create table %s: %w", s.tableName, err)
	}

	return nil
}

func NewCassandraStorage(session *gocql.Session, keyspace string, tableName string) *Storage {
	return &Storage{
		session:   session,
		keyspace:  keyspace,
		tableName: tableName,
	}
}
