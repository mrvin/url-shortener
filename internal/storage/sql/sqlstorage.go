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

	insertURL *sql.Stmt
	selectURL *sql.Stmt
	deleteURL *sql.Stmt

	selectURLs      *sql.Stmt
	selectTotalURLs *sql.Stmt

	existsAlias *sql.Stmt
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

func (s *Storage) CreateURL(ctx context.Context, username, url, alias string) error {
	if _, err := s.insertURL.ExecContext(ctx, url, alias, username); err != nil {
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

func (s *Storage) DeleteURL(ctx context.Context, username, alias string) error {
	res, err := s.deleteURL.ExecContext(ctx, username, alias)
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

func (s *Storage) GetURLs(ctx context.Context, username string, limit, offset uint64) ([]storage.URL, uint64, error) {
	urls := make([]storage.URL, 0)

	rows, err := s.selectURLs.QueryContext(ctx, username, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return urls, 0, nil
		}
		return nil, 0, fmt.Errorf("can't get rows urls: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url storage.URL
		err = rows.Scan(
			&url.URL,
			&url.Alias,
			&url.Count,
			&url.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("can't scan next row: %w", err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	var total uint64
	if err := s.selectTotalURLs.QueryRowContext(ctx, username).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("can't scan total urls: %w", err)
	}

	return urls, total, nil
}

func (s *Storage) CheckAlias(ctx context.Context, alias string) (bool, error) {
	var exists bool
	if err := s.existsAlias.QueryRowContext(ctx, alias).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists alias: %w", err)
	}

	return exists, nil
}

func (s *Storage) Close() error {
	s.insertUser.Close()
	s.selectUser.Close()

	s.insertURL.Close()
	s.selectURL.Close()
	s.deleteURL.Close()

	s.selectURLs.Close()
	s.selectTotalURLs.Close()

	s.existsAlias.Close()

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
		INSERT INTO urls (url, alias, username)
			VALUES($1, $2, $3)`
	s.insertURL, err = s.db.PrepareContext(ctx, sqlInsertURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert url", err)
	}
	const sqlSelectURL = `
		UPDATE urls
		SET count = count+1 
		WHERE alias = $1
		RETURNING url`

	s.selectURL, err = s.db.PrepareContext(ctx, sqlSelectURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select url", err)
	}
	const sqlDeleteURL = `
		DELETE FROM urls
		WHERE username = $1 AND alias = $2`
	s.deleteURL, err = s.db.PrepareContext(ctx, sqlDeleteURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete url", err)
	}
	const sqlSelectURLs = `
		SELECT 
			url,
			alias,
			count,
			created_at
		FROM urls
		WHERE username = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3`
	s.selectURLs, err = s.db.PrepareContext(ctx, sqlSelectURLs)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select urls", err)
	}

	const sqlSelectTotalURLs = `SELECT COUNT(alias) FROM urls WHERE username = $1`
	s.selectTotalURLs, err = s.db.PrepareContext(ctx, sqlSelectTotalURLs)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select total urls", err)
	}

	const sqlExistsAlias = `SELECT EXISTS ( SELECT 1 FROM urls WHERE alias = $1 )`
	s.existsAlias, err = s.db.PrepareContext(ctx, sqlExistsAlias)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "exists alias", err)
	}

	return nil
}
