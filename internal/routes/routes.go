package routes

import (
	"test-golang/internal/handlers"
	"test-golang/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	router *mux.Router,
	articleHandler *handlers.ArticleHandler,
	userHandler *handlers.UserHandler,
) {
	api := router.PathPrefix("/api/v1").Subrouter()

	SetupArticleRoutes(api, articleHandler)
	SetupUserRoutes(api, userHandler)

	// Global middleware
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
}
