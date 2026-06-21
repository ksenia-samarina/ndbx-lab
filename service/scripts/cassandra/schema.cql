CREATE KEYSPACE IF NOT EXISTS ${CASSANDRA_KEYSPACE}
            WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE ${CASSANDRA_KEYSPACE};

CREATE TABLE IF NOT EXISTS event_reactions
(
    event_id   text,
    created_by text,
    like_value tinyint,
    created_at timestamp,
    PRIMARY KEY (event_id, created_by)
);

CREATE TABLE IF NOT EXISTS event_reviews
(
    id         uuid,
    event_id   text,
    created_by text,
    rating     tinyint,
    comment    text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (event_id, id)
);

CREATE INDEX IF NOT EXISTS ON event_reactions (like_value);
CREATE INDEX IF NOT EXISTS ON event_reviews (rating);
