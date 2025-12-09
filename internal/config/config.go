package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/mrvin/url-shortener/internal/cache"
	"github.com/mrvin/url-shortener/internal/httpserver"
	"github.com/mrvin/url-shortener/internal/logger"
	sqlstorage "github.com/mrvin/url-shortener/internal/storage/sql"
)

type Config struct {
	DB     sqlstorage.Conf `yaml:"db"`
	Cache  cache.Conf      `yaml:"cache"`
	HTTP   httpserver.Conf `yaml:"http"`
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

	if host := os.Getenv("REDIS_HOST"); host != "" {
		c.Cache.Host = host
	} else {
		slog.Warn("Empty redis host")
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		c.Cache.Port = port
	} else {
		slog.Warn("Empty redis port")
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		c.Cache.Password = password
	} else {
		slog.Warn("Empty redis password")
	}
	if strName := os.Getenv("REDIS_DB"); strName != "" {
		if name, err := strconv.Atoi(strName); err != nil {
			slog.Warn("invalid redis db name: " + strName)
			c.Cache.Name = name
		}
	} else {
		slog.Warn("Empty redis db name")
	}

	if host := os.Getenv("HTTP_HOST"); host != "" {
		c.HTTP.Host = host
	} else {
		slog.Warn("Empty server http host")
	}
	if port := os.Getenv("HTTP_PORT"); port != "" {
		c.HTTP.Port = port
	} else {
		slog.Warn("Empty server http port")
	}
	if docFilePath := os.Getenv("DOC_FILEPATH"); docFilePath != "" {
		c.HTTP.DocFilePath = docFilePath
	} else {
		slog.Warn("Empty doc file path")
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
