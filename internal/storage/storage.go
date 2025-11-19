package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserExists   = errors.New("user exists")
	ErrUserNotFound = errors.New("user not found")

	ErrURLNotFound         = errors.New("url not found")
	ErrAliasExists         = errors.New("alias already exists")
	ErrURLAliasIsNotExists = errors.New("url alias is not exists")
)

//nolint:tagliatelle
type User struct {
	Name         string `json:"name"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

//nolint:tagliatelle
type URL struct {
	URL       string    `json:"url"`
	Alias     string    `json:"alias"`
	Count     string    `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)
}

type URLStorage interface {
	CreateURL(ctx context.Context, username, url, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, username, alias string) error
	GetURLs(ctx context.Context, username string, limit, offset uint64) ([]URL, uint64, error)
	CheckAlias(ctx context.Context, alias string) (bool, error)
}

type Storage interface {
	UserStorage
	URLStorage
}
