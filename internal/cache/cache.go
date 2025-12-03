package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrAliasNotFound = errors.New("alias not found")
)

const ttl = 30 * time.Minute

type Conf struct {
	Host     string
	Port     string
	Password string
	Name     int
}

type Cacher interface {
	GetURL(ctx context.Context, alias string) (string, error)
	SetURL(ctx context.Context, alias, url string) error
	DeleteURL(ctx context.Context, alias string) error
}

type Cache struct {
	conn *redis.Client
}

func New(ctx context.Context, conf *Conf) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{ //nolint:exhaustruct
		Addr:     net.JoinHostPort(conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Name,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache ping: %w", err)
	}

	return &Cache{rdb}, nil
}

func (c *Cache) Close() error {
	return c.conn.Close() //nolint:wrapcheck
}

func (c *Cache) GetURL(ctx context.Context, alias string) (string, error) {
	value, err := c.conn.Get(ctx, alias).Result()
	if errors.Is(err, redis.Nil) {
		slog.DebugContext(ctx, "Cache miss")
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting url from cache: %w", err)
	}

	return value, nil
}

func (c *Cache) SetURL(ctx context.Context, alias, url string) error {
	if err := c.conn.Set(ctx, alias, url, ttl).Err(); err != nil {
		return fmt.Errorf("setting url to cache: %w", err)
	}

	return nil
}

func (c *Cache) DeleteURL(ctx context.Context, alias string) error {
	count, err := c.conn.Del(ctx, alias).Result()
	if err != nil {
		return fmt.Errorf("deleting url from cache: %w", err)
	}
	if count == 0 {
		return ErrAliasNotFound
	}

	return nil
}
