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
	port := os.Getenv("APP_PORT")

	// handlers
	http.Handle("/health", check_health.New())

	log.Println("Listening http at port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalln("Cannot listen http", err)
	}
}
