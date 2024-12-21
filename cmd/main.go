package main

import (
	"log"
	"net/http"
	"giproxy/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	router := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}