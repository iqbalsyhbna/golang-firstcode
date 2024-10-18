package routes

import (
	"golang-firstcode/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupArticleRoutes(router *mux.Router, handler *handlers.ArticleHandler) {
	articles := router.PathPrefix("/articles").Subrouter()
	articles.HandleFunc("", handler.GetArticles).Methods(http.MethodGet)
	articles.HandleFunc("/{id}", handler.GetArticle).Methods(http.MethodGet)
	articles.HandleFunc("", handler.CreateArticle).Methods(http.MethodPost)
	articles.HandleFunc("/{id}", handler.UpdateArticle).Methods(http.MethodPut)
	articles.HandleFunc("/{id}", handler.DeleteArticle).Methods(http.MethodDelete)
}
