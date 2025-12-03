package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
)

type DBURLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
	CountIncrement(alias string) error
}

type CacheURLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
	SetURL(ctx context.Context, alias, url string) error
}

func NewRedirect(st DBURLGetter, cache CacheURLGetter) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		alias := req.PathValue("alias")
		ctx := logger.WithAlias(req.Context(), alias)
		msg := "Redirect"

		url, err := cache.GetURL(ctx, alias)
		if err != nil {
			err = fmt.Errorf("getting url from cache: %w", err)
			slog.WarnContext(ctx, msg, slog.String("warn", err.Error()))
		}
		if url == "" {
			url, err = st.GetURL(ctx, alias)
			if err != nil {
				err = fmt.Errorf("getting url from storage: %w", err)
				if errors.Is(err, storage.ErrAliasNotFound) {
					return ctx, http.StatusNotFound, err
				}
				return ctx, http.StatusInternalServerError, err
			}
			ctx = logger.WithURL(ctx, url)
			if err := cache.SetURL(ctx, alias, url); err != nil {
				slog.WarnContext(ctx, msg, slog.String("warn", err.Error()))
			}
		} else {
			ctx = logger.WithURL(ctx, url)
			go func() {
				if err := st.CountIncrement(alias); err != nil {
					slog.WarnContext(ctx, msg, slog.String("warn", err.Error()))
				}
			}()
		}

		// redirect to found url
		http.Redirect(res, req, url, http.StatusFound)

		return ctx, http.StatusFound, nil
	}
}
