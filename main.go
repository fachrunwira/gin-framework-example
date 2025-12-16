package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fachrunwira/gin-example/database"
	"github.com/fachrunwira/gin-example/lib/env"
	"github.com/fachrunwira/gin-example/lib/logger"
	"github.com/fachrunwira/gin-example/middlewares"
	"github.com/fachrunwira/gin-example/middlewares/ratelimit"
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

	appLogger = logger.New("./storage/logs/app_log.log", slog.LevelError)
	rateLimitLogger = logger.New("./storage/logs/rate_limit.log", slog.LevelInfo)

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

	rl := ratelimit.New(rate.Every(45*time.Second), 100, time.Minute, 5*time.Minute, rateLimitLogger)
	defer rl.Stop()

	g.Use(rl.Middleware())
	g.Use(middlewares.SetLog("./storage/logs/http.log"))
	g.Use(middlewares.InjectDB())

	routes.RegisterRoutes(g)

	port := env.GetEnv("APP_PORT", "8080")

	g.Run(":" + port)

	// Gracefull shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: g,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exit")
}
