package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rashaev/todo-app/config"
	"github.com/rashaev/todo-app/internal/handler"
	"github.com/rashaev/todo-app/internal/repository/database"
	"github.com/rashaev/todo-app/internal/usecase"
	"github.com/rashaev/todo-app/pkg/postgres"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// Init DB
	db, err := postgres.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		log.Fatalf("unable to create db connection: %v\n", err)
	}
	defer db.Close()

	// Init repo
	todoRepo := database.NewTodoRepository(db)

	// Init UseCase
	todoUseCase := usecase.NewTodoUseCase(todoRepo)

	// Init handlers
	todoHandlers := handler.NewTodoHandlers(todoUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/todos", todoHandlers.GetAllTodos).Methods("GET")
	router.HandleFunc("/todos", todoHandlers.CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandlers.GetTodoByID).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandlers.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandlers.DeleteTodo).Methods("DELETE")

	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: router,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	// Run server in goroutine
	go func() {
		log.Printf("Server is running on %s\n", cfg.ListenAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Waiting for SIGTERM or SIGKILL signals
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
