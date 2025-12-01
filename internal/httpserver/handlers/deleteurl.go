package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
)

type DBURLDeleter interface {
	DeleteURL(ctx context.Context, username, alias string) error
}

type CacheURLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func NewDeleteURL(st DBURLDeleter, cache CacheURLDeleter) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		alias := req.PathValue("alias")
		ctx := logger.WithAlias(req.Context(), alias)
		msg := "Delete url"

		username, err := logger.GetUsernameFromCtx(ctx)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("get user name from ctx: %w", err)
		}

		if err := st.DeleteURL(ctx, username, alias); err != nil {
			err = fmt.Errorf("deleting url from storage: %w", err)
			if errors.Is(err, storage.ErrAliasNotFound) {
				return ctx, http.StatusNotFound, err
			}
			return ctx, http.StatusInternalServerError, err
		}
		if err := cache.DeleteURL(ctx, alias); err != nil {
			err = fmt.Errorf("deleting url from cache: %w", err)
			slog.WarnContext(ctx, msg, slog.String("warn", err.Error()))
		}

		httpresponse.WriteOK(res, http.StatusOK)

		return ctx, http.StatusOK, nil
	}
}
