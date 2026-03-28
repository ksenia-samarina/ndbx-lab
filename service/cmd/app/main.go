package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/domains/session"
	"samarina/ndbx/internal/handlers/health"
	session2 "samarina/ndbx/internal/handlers/session"
	"samarina/ndbx/internal/repository/session/redis"
	"strconv"
	"time"

	redisdb "github.com/redis/go-redis/v9"
)

func main() {
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

	client := redisdb.NewClient(&redisdb.Options{
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
	}(client)

	storage := redis.NewSessionStorage(client)
	domain := session.NewDomain(storage)

	// handlers
	http.Handle("/health", health.New(ttl))
	http.Handle("/session", session2.New(domain, ttl))

	log.Printf("Listening http at addr: %s:%s", appHost, appPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", appHost, appPort), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
