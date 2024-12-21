package main

import (
	"log"
	"net/http"
	"giproxy/internal/routes"
	"github.com/joho/godotenv"
	// "os"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	// log.Printf("GIT_TOKEN: %s", os.Getenv("GIT_TOKEN")) // Uncomment this line to print the GIT_TOKEN

	router := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}