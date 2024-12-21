package main

import (
	"log"
	"net/http"
	"giproxy/internal/routes"
)

func main() {
	router := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}