package routes

import (
	"giproxy/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/api/git/").HandlerFunc(handlers.CommonHandler)

	return router
}