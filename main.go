package main

import (
	"log"
	"net/http"
	"test-golang/internal/config"
	"test-golang/internal/handlers"
	"test-golang/internal/routes"
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
	userService := service.NewUserService(db)

	// Initialize handlers with services
	articleHandler := handlers.NewArticleHandler(articleService)
	userHandler := handlers.NewUserHandler(userService)

	routes.SetupRoutes(app.Router, articleHandler, userHandler)
	return app
}

func (app *App) Run(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, app.Router)
}

func main() {
	app := NewApp()
	log.Fatal(app.Run(":8000"))
}
