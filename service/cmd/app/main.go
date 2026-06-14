package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/domains/auth/login"
	"samarina/ndbx/internal/domains/auth/logout"
	"samarina/ndbx/internal/domains/events"
	"samarina/ndbx/internal/domains/reactions"
	"samarina/ndbx/internal/domains/session"
	"samarina/ndbx/internal/domains/users"
	"samarina/ndbx/internal/domains/validator"
	login2 "samarina/ndbx/internal/handlers/auth/login"
	logout2 "samarina/ndbx/internal/handlers/auth/logout"
	events2 "samarina/ndbx/internal/handlers/events"
	"samarina/ndbx/internal/handlers/health"
	reactions2 "samarina/ndbx/internal/handlers/reactions"
	session2 "samarina/ndbx/internal/handlers/session"
	users2 "samarina/ndbx/internal/handlers/users"
	storagecassandra "samarina/ndbx/internal/repository/cassandra/reactions"
	storageevents "samarina/ndbx/internal/repository/mongo/events"
	storageusers "samarina/ndbx/internal/repository/mongo/users"
	"samarina/ndbx/internal/repository/redis"
	"strconv"
	"strings"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	redisdb "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")

	sec, err := strconv.Atoi(os.Getenv("APP_USER_SESSION_TTL"))
	if err != nil {
		log.Fatal(err)
	}
	ttl := time.Duration(sec) * time.Second

	likeTTLSecStr := os.Getenv("APP_LIKE_TTL")
	if likeTTLSecStr == "" {
		likeTTLSecStr = "60"
	}
	likeTTLSec, err := strconv.Atoi(likeTTLSecStr)
	if err != nil {
		log.Fatal(err)
	}
	likeTTL := time.Duration(likeTTLSec) * time.Second

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}

	mongoHost := os.Getenv("MONGODB_HOST")
	mongoPort := os.Getenv("MONGODB_PORT")
	mongoDB := os.Getenv("MONGODB_DATABASE")

	cassandraHostsStr := os.Getenv("CASSANDRA_HOSTS")
	cassandraPortStr := os.Getenv("CASSANDRA_PORT")
	cassandraUsername := os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword := os.Getenv("CASSANDRA_PASSWORD")
	cassandraKeyspace := os.Getenv("CASSANDRA_KEYSPACE")
	cassandraConsistency := os.Getenv("CASSANDRA_CONSISTENCY")

	redisClient := redisdb.NewClient(&redisdb.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,

		PoolSize:     10,
		MinIdleConns: 3,
	})
	defer func(client *redisdb.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(redisClient)

	mongoURI := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(3).
		SetMaxConnIdleTime(time.Minute * 5)

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func(client *mongo.Client) {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(client)
		}
	}(mongoClient)

	db := mongoClient.Database(mongoDB)

	hosts := strings.Split(cassandraHostsStr, ",")
	cluster := gocql.NewCluster(hosts...)
	if cassandraPortStr != "" {
		port, err := strconv.Atoi(cassandraPortStr)
		if err == nil {
			cluster.Port = port
		}
	}

	if cassandraUsername != "" && cassandraPassword != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cassandraUsername,
			Password: cassandraPassword,
		}
	}

	cluster.Consistency = gocql.One
	if strings.ToUpper(cassandraConsistency) == "QUORUM" {
		cluster.Consistency = gocql.Quorum
	}

	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	gocqlSession, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to connect to Cassandra cluster: %v", err)
	}
	defer gocqlSession.Close()

	redisStorage := redis.NewStorage(redisClient)
	sessionDomain := session.NewDomain(redisStorage)

	mongodbUsersStorage := storageusers.NewStorage(db, "users")
	mongodbEventsStorage := storageevents.NewStorage(db, "events", "users")

	cassandraStorage := storagecassandra.NewCassandraStorage(gocqlSession, cassandraKeyspace, "event_reactions")

	initSchemaCtx, cancelSchema := context.WithTimeout(ctx, 10*time.Second)
	err = cassandraStorage.InitSchema(initSchemaCtx, cassandraKeyspace)
	cancelSchema()
	if err != nil {
		log.Fatalf("failed to init cassandra schema: %v", err)
	}

	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	err = mongodbUsersStorage.CreateIndexes(ctx, userIndexes)
	if err != nil {
		log.Fatal(err)
	}

	eventIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "title", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "title", Value: 1},
				{Key: "created_by", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "created_by", Value: 1}},
		},
	}
	err = mongodbEventsStorage.CreateIndexes(ctx, eventIndexes)
	if err != nil {
		log.Fatal(err)
	}

	usersDomain := users.NewDomain(redisStorage, mongodbUsersStorage, mongodbEventsStorage)
	eventsDomain := events.NewDomain(redisStorage, mongodbEventsStorage)
	validatorDomain := validator.NewDomain()

	reactionsDomain := reactions.NewDomain(cassandraStorage, mongodbEventsStorage, redisStorage, likeTTL)

	loginDomain := login.NewDomain(redisStorage, mongodbUsersStorage)
	logoutDomain := logout.NewDomain(redisStorage)

	// handlers
	http.Handle("/health", health.New(ttl))
	http.Handle("/session", session2.New(sessionDomain, ttl))

	usersHandler := users2.New(usersDomain, validatorDomain, reactionsDomain, ttl)
	eventsHandler := events2.New(eventsDomain, reactionsDomain, validatorDomain, ttl)
	reactionsHandler := reactions2.New(eventsDomain, reactionsDomain, ttl)

	http.HandleFunc("/users", usersHandler.RegisterOrGetUsers)
	http.HandleFunc("/users/{id}", usersHandler.GetUserByID)
	http.HandleFunc("/users/{id}/events", usersHandler.GetUserEventsByUserID)

	http.HandleFunc("/events", eventsHandler.RegisterOrGetEvents)
	http.HandleFunc("/events/{id}", eventsHandler.GetOrEditEventData)

	http.Handle("/auth/login", login2.New(loginDomain, ttl))
	http.Handle("/auth/logout", logout2.New(logoutDomain, ttl))

	http.HandleFunc("/events/{event_id}/like", reactionsHandler.SetLike)
	http.HandleFunc("/events/{event_id}/dislike", reactionsHandler.SetDislike)

	log.Printf("Listening http at addr: %s:%s", appHost, appPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", appHost, appPort), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
