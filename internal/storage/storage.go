package storage

import (
	"context"
	"errors"
)

var (
	ErrUserExists   = errors.New("user exists")
	ErrUserNotFound = errors.New("user not found")

	ErrURLNotFound         = errors.New("url not found")
	ErrURLExists           = errors.New("url exists")
	ErrURLAliasIsNotExists = errors.New("url alias is not exists")
)

//nolint:tagliatelle
type User struct {
	Name         string `json:"name"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)
}

type URLStorage interface {
	CreateURL(ctx context.Context, userName, urlToSave, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, userName, alias string) error
	GetCountURL(ctx context.Context, userName, alias string) (int64, error)
}

type Storage interface {
	UserStorage
	URLStorage
}
