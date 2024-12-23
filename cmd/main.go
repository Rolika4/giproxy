package main

import (
	"giproxy/internal/routes"
	"log"
	"net/http"
	// "github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	//     log.Fatalf("Error loading .env file: %v", err)
	// }

	router := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
