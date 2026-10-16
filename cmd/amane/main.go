package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"idac-amane-bot/internal/adapter/database"
	discord_handler "idac-amane-bot/internal/adapter/discord"
	"idac-amane-bot/internal/config"
	"idac-amane-bot/pkg/logger"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// Init structured logger (level from APP_ENV: "dev" => debug, else info).
	logger.Init(os.Getenv("APP_ENV"))

	// 2. Connect to DB
	dbConn, err := database.NewPostgresDBConn(cfg.DatabaseCfg)
	if err != nil {
		log.Fatal("Postgres database connection error: ", err)
	}
	// dbConn is wired and ready for repositories (internal/adapter/repository).
	defer dbConn.Close()

	// 3. Init Repositories
	// (none yet — add repositories here as features grow)

	// 4. Init Discord Session
	discordSession, err := discord_handler.NewDiscordSession(&cfg.DiscordCfg)
	if err != nil {
		log.Fatal("Discord session creation error:", err)
	}
	if err = discordSession.OpenSession(); err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer discordSession.CloseSession()

	// 5. Init services
	// (none yet — add domain services here as features grow)

	cmdHandler := discord_handler.NewHandler(
		discordSession.Session,
		cfg.DiscordCfg.BotOwnerID,
	)
	// 6. Register Commands & Event Handlers
	if err := cmdHandler.RegisterCommands(); err != nil {
		logger.Error(context.Background(), "failed to register commands", err)
	}

	// 7. Health server
	// The bot is outbound-only (Discord gateway, HTTP, Postgres) and never accepts
	// inbound traffic. Web-service hosts like Render port-scan the container and
	// route health checks to an open port, so we expose a tiny HTTP endpoint purely
	// to satisfy that. Bind 0.0.0.0 (not localhost) or the scan won't detect it.
	// Deployed as a background worker, this is simply unused.
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Render's default
	}
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	healthSrv := &http.Server{Addr: "0.0.0.0:" + port, Handler: healthMux}
	go func() {
		logger.Info(context.Background(), fmt.Sprintf("health server listening on :%s", port))
		if err := healthSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "health server error", err)
		}
	}()

	// 8. Init graceful shutdown
	logger.Info(context.Background(), "bot is running, press CTRL-C to exit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info(context.Background(), "gracefully shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := healthSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error(context.Background(), "health server shutdown error", err)
	}
}
