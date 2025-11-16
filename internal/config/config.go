package config

import (
	"log/slog"
	"os"

	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver"
	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/url-shortener/internal/storage/sql"
)

type Config struct {
	HTTP   httpserver.Conf `yaml:"http"`
	DB     sqlstorage.Conf `yaml:"db"`
	Logger logger.Conf     `yaml:"logger"`
}

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		c.DB.Host = host
	} else {
		slog.Warn("Empty postgres host")
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		c.DB.Port = port
	} else {
		slog.Warn("Empty postgres port")
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		c.DB.User = user
	} else {
		slog.Warn("Empty postgres user")
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		c.DB.Password = password
	} else {
		slog.Warn("Empty postgres password")
	}
	if name := os.Getenv("POSTGRES_DB"); name != "" {
		c.DB.Name = name
	} else {
		slog.Warn("Empty postgres db name")
	}

	if addr := os.Getenv("SERVER_HTTP_ADDR"); addr != "" {
		c.HTTP.Addr = addr
	} else {
		slog.Warn("Empty server http addr")
	}

	if logFilePath := os.Getenv("LOGGER_FILEPATH"); logFilePath != "" {
		c.Logger.FilePath = logFilePath
	} else {
		slog.Warn("Empty log file path")
	}
	if logLevel := os.Getenv("LOGGER_LEVEL"); logLevel != "" {
		c.Logger.Level = logLevel
	} else {
		slog.Warn("Empty log level")
	}
}
