package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/domains/session"
	"samarina/ndbx/internal/handlers/check_health"
	"samarina/ndbx/internal/handlers/update_session"
	"samarina/ndbx/internal/repository/redis"
	"strconv"

	redisdb "github.com/redis/go-redis/v9"
)

func main() {
	// read os.Getenv
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")

	redisHost := os.Getenv("REDIS_CONTAINER_NAME")
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
	domain := session.New(storage)

	// handlers
	http.Handle("/health", check_health.New())
	http.Handle("/session", update_session.New(domain))

	log.Printf("Listening http at addr: %s:%s", appHost, appPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", appHost, appPort), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
