package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fachrunwira/gin-example/app/middlewares"
	"github.com/fachrunwira/gin-example/database"
	"github.com/fachrunwira/gin-example/lib"
	"github.com/fachrunwira/gin-example/lib/env"
	"github.com/fachrunwira/gin-example/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var (
	appLogger       *slog.Logger
	rateLimitLogger *slog.Logger
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatalln(errEnv)
	}

	appLogger = lib.NewLogger("./storage/logs/app_log.log", slog.LevelError)
	rateLimitLogger = lib.NewLogger("./storage/logs/rate_limit.log", slog.LevelInfo)

	dbOptions := &database.DatabaseOptions{
		MaxOpenConnection:     25,
		MaxIdleConnection:     25,
		MaxConnectionLifetime: 2 * time.Minute,
	}
	if errDB := database.Init(dbOptions); errDB != nil {
		appLogger.Error("cannot connect to", "db", errDB)
		return
	}
	defer database.Close()

	g := gin.Default()

	rl := middlewares.NewRateLimit(rate.Every(45*time.Second), 100, time.Minute, 5*time.Minute, []string{}, rateLimitLogger)
	defer rl.Stop()

	g.Use(rl.Middleware())
	g.Use(middlewares.SetLog("./storage/logs/http.log"))

	gin.SetMode(gin.TestMode)

	routes.RegisterRoutes(g)

	port := env.Get("APP_PORT", "8080")

	g.Run(":" + port)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        g,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	go startServer(server, appLogger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	gracefullShutdown(context.Background(), 5*time.Second, server, appLogger)
}

func startServer(server *http.Server, logger *slog.Logger) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
	}
}

func gracefullShutdown(parent context.Context, timeout time.Duration, server *http.Server, logger *slog.Logger) {
	logger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exit")
}
