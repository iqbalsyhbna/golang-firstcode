package main

import (
	"log"
	"net/http"
	"test-golang/handlers"
	"test-golang/middleware"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/articles", handlers.GetArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", handlers.GetArticle).Methods("GET")
	router.HandleFunc("/articles", handlers.CreateArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", handlers.UpdateArticle).Methods("PUT")
	router.HandleFunc("/articles/{id}", handlers.DeleteArticle).Methods("DELETE")

	// Apply CORS middleware
	router.Use(middleware.CORS)

	log.Fatal(http.ListenAndServe(":8000", router))
}
