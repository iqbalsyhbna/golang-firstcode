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
	config.InitDBs()

	app := &App{
		Router: mux.NewRouter(),
	}

	golangDB := config.GetDB("golang_db")

	if golangDB == nil {
		log.Fatal("Failed to initialize one or more databases")
	}

	// Initialize services with database connection
	articleService := service.NewArticleService(golangDB)
	userService := service.NewUserService(golangDB)

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

	defer func() {
		for dbName, db := range config.DBMap {
			if err := db.Close(); err != nil {
				log.Printf("Error closing %s database connection: %v", dbName, err)
			}
		}
	}()

	log.Fatal(app.Run(":8000"))
}
