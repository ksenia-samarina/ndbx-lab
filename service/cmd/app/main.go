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
	"samarina/ndbx/internal/domains/session"
	"samarina/ndbx/internal/domains/users"
	login2 "samarina/ndbx/internal/handlers/auth/login"
	logout2 "samarina/ndbx/internal/handlers/auth/logout"
	events2 "samarina/ndbx/internal/handlers/events"
	"samarina/ndbx/internal/handlers/health"
	session2 "samarina/ndbx/internal/handlers/session"
	users2 "samarina/ndbx/internal/handlers/users"
	mongodb "samarina/ndbx/internal/repository/mongo"
	"samarina/ndbx/internal/repository/redis"
	"strconv"
	"time"

	redisdb "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// read os.Getenv
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")

	sec, err := strconv.Atoi(os.Getenv("APP_USER_SESSION_TTL"))
	if err != nil {
		log.Fatal(err)
	}
	ttl := time.Duration(sec) * time.Second

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
	//mongoUser := os.Getenv("MONGODB_USER")
	//mongoPassword := os.Getenv("MONGODB_PASSWORD")

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
			log.Fatal(err)
		}
	}(mongoClient)

	db := mongoClient.Database(mongoDB)

	redisStorage := redis.NewStorage(redisClient)
	sessionDomain := session.NewDomain(redisStorage)

	mongodbUsersStorage := mongodb.NewStorage(db, "users")
	mongodbEventsStorage := mongodb.NewStorage(db, "events")

	// indexes
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
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: options.Index().SetUnique(true),
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

	usersDomain := users.NewDomain(redisStorage, mongodbUsersStorage)
	eventsDomain := events.NewDomain(redisStorage, mongodbEventsStorage)

	loginDomain := login.NewDomain(redisStorage, mongodbUsersStorage)
	logoutDomain := logout.NewDomain(redisStorage)

	// handlers
	http.Handle("/health", health.New(ttl))
	http.Handle("/session", session2.New(sessionDomain, ttl))
	http.Handle("/users", users2.New(usersDomain, ttl))
	http.Handle("/events", events2.New(eventsDomain, ttl))
	http.Handle("/auth/login", login2.New(loginDomain, ttl))
	http.Handle("/auth/logout", logout2.New(logoutDomain, ttl))

	log.Printf("Listening http at addr: %s:%s", appHost, appPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", appHost, appPort), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
