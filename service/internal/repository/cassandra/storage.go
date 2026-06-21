package cassandra

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Storage struct {
	session         *gocql.Session
	keyspace        string
	likeTableName   string
	reviewTableName string
}

func NewCassandraStorage(session *gocql.Session, keyspace string, likeTableName string, reviewTableName string) *Storage {
	return &Storage{
		session:         session,
		keyspace:        keyspace,
		likeTableName:   likeTableName,
		reviewTableName: reviewTableName,
	}
}
