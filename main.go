package main

import (
	"aulway/internal/database/postgres"
	"aulway/internal/firebase"
	xtransport "aulway/internal/transport/http"
	"aulway/internal/utils/config"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/oklog/run"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("error getting config:", "error", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("database connection failed:", "error", err.Error())
		return
	}

	slog.Info("Database connection success")

	sqlDB, err := database.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm.DB:", "error", err.Error())
		return
	}

	authClient, err := firebase.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	router := xtransport.NewRouter(cfg, database, authClient).Build()

	var g run.Group
	{
		g.Add(func() error {
			return router.Start(":" + cfg.Port)
		}, func(err error) {
			ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				cancel()
				if shutdownErr := router.Shutdown(ctxShutDown); shutdownErr != nil {
					slog.Error("error shutting down router", "error", shutdownErr.Error())
				}

				if dbErr := sqlDB.Close(); dbErr != nil {
					slog.Error("error closing database", "error", dbErr.Error())
				}
			}()
			//e := router.Shutdown(ctxShutDown)
			//if e != nil {
			//	os.Exit(1)
			//}
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			cs := make(chan os.Signal, 1)
			signal.Notify(cs, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-cs:
				return fmt.Errorf("recieved signal %+v", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(err error) {
			close(cancelInterrupt)

			if dbErr := sqlDB.Close(); dbErr != nil {
				slog.Error("error closing database", "error", dbErr.Error())
			}
		})
	}
	log.Warn("exit", g.Run())
}
