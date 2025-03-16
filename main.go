package main

import (
	"aulway/internal/database/postgres"
	"aulway/internal/database/redis"
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

// @title Aulway API
// @version 1.0
// @description API documentation for Aulway.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use "Bearer {your-firebase-token}"
// @BasePath /
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("error getting config:", "error", err.Error())
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("database connection failed:", "error", err.Error())
		panic(err)
	}

	slog.Info("Database connection success")

	sqlDB, err := database.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm.DB:", "error", err.Error())
		panic(err)
	}

	redis, err := redis.Connect(ctx, cfg.Redis)
	if err != nil {
		slog.Error("redis connection failed:", "error", err.Error())
		panic(err)
	}

	router := xtransport.NewRouter(cfg, database, redis).Build()

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
