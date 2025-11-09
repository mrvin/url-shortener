package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver"
	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/url-shortener/internal/storage/sql"
)

//nolint:tagliatelle
type Config struct {
	DefaultAliasLength int             `yaml:"default_alias_length"`
	HTTP               httpserver.Conf `yaml:"http"`
	DB                 sqlstorage.Conf `yaml:"db"`
	Logger             logger.Conf     `yaml:"logger"`
}

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if defaultAliasLength := os.Getenv("DEFAULT_ALIAS_LENGTH"); defaultAliasLength != "" {
		defaultAliasLengthInt, err := strconv.Atoi(defaultAliasLength)
		if err != nil {
			slog.Warn("Invalid default alias length " + defaultAliasLength)
		}
		c.DefaultAliasLength = defaultAliasLengthInt
	} else {
		slog.Warn("Empty default alias length")
	}
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
