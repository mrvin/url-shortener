package main

import (
	"context"
	stdlog "log"
	"log/slog"

	"github.com/mrvin/tasks-go/url-shortener/internal/config"
	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver"
	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/url-shortener/internal/storage/sql"
)

func main() {
	// init config
	var conf config.Config
	conf.LoadFromEnv()

	// init logger
	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v\n", err)
		return
	}
	slog.Info("Init logger", slog.String("level", conf.Logger.Level))
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()

	// init storage
	ctx := context.Background()
	st, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		slog.Error("Failed to init storage: " + err.Error())
		return
	}
	slog.Info("Connected to database")
	defer func() {
		if err := st.Close(); err != nil {
			slog.Error("Failed to close storage: " + err.Error())
		} else {
			slog.Info("Closing the database connection")
		}
	}()

	// Start server
	server := httpserver.New(&conf.HTTP, conf.DefaultAliasLength, st)

	server.Run(ctx)
}
