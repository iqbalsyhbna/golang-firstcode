package routes

import (
	"net/http"
	"test-golang/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
	users := router.PathPrefix("/users").Subrouter()
	users.HandleFunc("", handler.GetUsers).Methods(http.MethodGet)
	users.HandleFunc("/{id}", handler.GetUser).Methods(http.MethodGet)
	users.HandleFunc("", handler.CreateUser).Methods(http.MethodPost)
	users.HandleFunc("/{id}", handler.UpdateUser).Methods(http.MethodPut)
	users.HandleFunc("/{id}", handler.DeleteUser).Methods(http.MethodDelete)
}
