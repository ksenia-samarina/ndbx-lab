package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/handlers/check_health"
)

func main() {
	// read os.Getenv
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	// handlers
	http.Handle("/health", check_health.New())

	log.Printf("Listening http at addr: %s:%s", host, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
