package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/config"
	"github.com/0DayMonxrch/project-management-system/internal/handler"
	"github.com/0DayMonxrch/project-management-system/internal/middleware"
	"github.com/0DayMonxrch/project-management-system/internal/repository"
	"github.com/0DayMonxrch/project-management-system/internal/service"
	"github.com/0DayMonxrch/project-management-system/migrations"
	"github.com/0DayMonxrch/project-management-system/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	// Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Logger
	log := logger.New(cfg.App.Env)

	// MongoDB
	client, err := repository.NewMongoClient(cfg.DB.URI)
	if err != nil {
		log.Error("failed to connect to mongodb", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Error("failed to disconnect mongodb", "error", err)
		}
	}()

	db := client.Database(cfg.DB.Name)
	log.Info("connected to mongodb", "database", cfg.DB.Name)

	// Run migrations
	if err := migrations.RunIndexes(db, log); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	noteRepo := repository.NewNoteRepository(db)

	// Services
	emailSvc := service.NewEmailService(cfg.SMTP)
	authSvc := service.NewAuthService(userRepo, emailSvc, cfg.JWT)
	projectSvc := service.NewProjectService(projectRepo, userRepo)
	taskSvc := service.NewTaskService(taskRepo, projectRepo)
	noteSvc := service.NewNoteService(noteRepo, projectRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	noteHandler := handler.NewNoteHandler(noteSvc)

	// Router
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authHandler, projectHandler, taskHandler, noteHandler, cfg.JWT.AccessSecret)

	// Global middleware chain: recovery → logger → router
	chain := middleware.Recovery(log)(middleware.Logger(log)(mux))

	// Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      chain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server starting", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
}
