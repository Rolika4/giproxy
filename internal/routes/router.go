package routes

import (
	"giproxy/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/api/git/").HandlerFunc(handlers.CommonHandler)
	router.PathPrefix("/widgets/").HandlerFunc(handlers.WidgetHandler)

	return router
}