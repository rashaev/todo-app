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
	"github.com/rashaev/todo-app/pkg/logger"
	"github.com/rashaev/todo-app/pkg/postgres"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Init logger
	logger := logger.NewSlogLogger(cfg.Logging.Level)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// Init DB
	db, err := postgres.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		logger.Error("unable to create db connection", "error", err)
	}
	defer db.Close()

	// Init repo
	todoRepo := database.NewTodoRepository(db)

	// Init UseCase
	todoUseCase := usecase.NewTodoUseCase(todoRepo)

	// Init handlers
	todoHandlers := handler.NewTodoHandlers(todoUseCase, logger)

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
		logger.Info("Server is running", "address", cfg.ListenAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("ListenAndServe error", "error", err)
		}
	}()

	// Waiting for SIGTERM or SIGKILL signals
	<-ctx.Done()
	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Info("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
