// cmd/api/main.go
package main

import (
	"log"
	"net/http"
	"test-golang/internal/config"
	"test-golang/internal/handlers"
	"test-golang/internal/middleware"
	"test-golang/internal/service"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func NewApp() *App {
	// Initialize database
	db := config.InitDB()

	app := &App{
		Router: mux.NewRouter(),
	}
	
	// Initialize services with database connection
	articleService := service.NewArticleService(db)
	
	// Initialize handlers with services
	articleHandler := handlers.NewArticleHandler(articleService)
	
	app.initializeRoutes(articleHandler)
	return app
}

func (app *App) initializeRoutes(articleHandler *handlers.ArticleHandler) {
	api := app.Router.PathPrefix("/api/v1").Subrouter()

	articles := api.PathPrefix("/articles").Subrouter()
	articles.HandleFunc("", articleHandler.GetArticles).Methods(http.MethodGet)
	articles.HandleFunc("/{id}", articleHandler.GetArticle).Methods(http.MethodGet)
	articles.HandleFunc("", articleHandler.CreateArticle).Methods(http.MethodPost)
	articles.HandleFunc("/{id}", articleHandler.UpdateArticle).Methods(http.MethodPut)
	articles.HandleFunc("/{id}", articleHandler.DeleteArticle).Methods(http.MethodDelete)

	app.Router.Use(middleware.Logging)
	app.Router.Use(middleware.CORS)
}

func (app *App) Run(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, app.Router)
}

func main() {
	app := NewApp()
	log.Fatal(app.Run(":8000"))
}