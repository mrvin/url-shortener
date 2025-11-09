package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

const maxOpenConns = 25
const maxIdleConns = 25
const connMaxLifetime = 5 * time.Minute

type Conf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertUser *sql.Stmt
	selectUser *sql.Stmt

	insertURL    *sql.Stmt
	selectURLOld *sql.Stmt
	selectURL    *sql.Stmt
	deleteURL    *sql.Stmt

	selectCountURL *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.connect(ctx); err != nil {
		return nil, err
	}
	if err := st.prepareQuery(ctx); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	if _, err := s.insertUser.ExecContext(ctx, user.Name, user.HashPassword, user.Role); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				return storage.ErrUserExists
			}
		}

		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, name string) (*storage.User, error) {
	var user storage.User

	if err := s.selectUser.QueryRowContext(ctx, name).Scan(&user.HashPassword, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't scan user with name: %s: %w", name, err)
	}
	user.Name = name

	return &user, nil
}

func (s *Storage) CreateURL(ctx context.Context, userName, urlToSave, alias string) error {
	if _, err := s.insertURL.ExecContext(ctx, urlToSave, alias, userName); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				return storage.ErrURLExists
			}
		}
		return fmt.Errorf("insert url: %w", err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	var url string

	if err := s.selectURL.QueryRowContext(ctx, alias).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("can't scan URL with alias: %s: %w", alias, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, userName, alias string) error {
	res, err := s.deleteURL.ExecContext(ctx, userName, alias)
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	if count != 1 {
		return storage.ErrURLAliasIsNotExists
	}

	return nil
}

func (s *Storage) GetCountURL(ctx context.Context, userName, alias string) (int64, error) {
	var count int64

	if err := s.selectCountURL.QueryRowContext(ctx, userName, alias).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrURLNotFound
		}
		return 0, fmt.Errorf("can't scan count for alias: %w", err)
	}

	return count, nil
}

func (s *Storage) Close() error {
	s.insertUser.Close()
	s.selectUser.Close()

	s.insertURL.Close()
	s.selectURLOld.Close()
	s.selectURL.Close()
	s.deleteURL.Close()

	s.selectCountURL.Close()

	return s.db.Close() //nolint:wrapcheck
}

func (s *Storage) connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	// Users query.
	sqlInsertUser := `
		INSERT INTO users (
			name,
			hash_password,
			role
		)
		VALUES ($1, $2, $3)`
	s.insertUser, err = s.db.PrepareContext(ctx, sqlInsertUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert user", err)
	}
	sqlGetUser := `
		SELECT hash_password,
			role
		FROM users
		WHERE name = $1`
	s.selectUser, err = s.db.PrepareContext(ctx, sqlGetUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select user", err)
	}

	// URL query.
	const sqlInsertURL = `
		INSERT INTO url(url, alias, user_name)
			VALUES($1, $2, $3)`
	s.insertURL, err = s.db.PrepareContext(ctx, sqlInsertURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert url", err)
	}
	const sqlSelectURLOld = `
		UPDATE 
			url
		SET count = count+1 
		WHERE alias = $1
		RETURNING url`

	s.selectURLOld, err = s.db.PrepareContext(ctx, sqlSelectURLOld)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select url", err)
	}
	const sqlSelectURL = "SELECT get_url($1)"
	s.selectURL, err = s.db.PrepareContext(ctx, sqlSelectURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select get_url func", err)
	}
	const sqlDeleteURL = `
		DELETE FROM url
		WHERE user_name = $1 AND alias = $2`
	s.deleteURL, err = s.db.PrepareContext(ctx, sqlDeleteURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete url", err)
	}
	const sqlSelectCountURL = `
		SELECT
			count
		FROM url
		WHERE user_name = $1 AND alias = $2`
	s.selectCountURL, err = s.db.PrepareContext(ctx, sqlSelectCountURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select count", err)
	}

	return nil
}
