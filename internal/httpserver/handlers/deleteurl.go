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

func NewDeleteURL(st DBURLDeleter, cache CacheURLDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")
		ctx := logger.WithAlias(req.Context(), alias)

		username, err := logger.GetUsernameFromCtx(ctx)
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(ctx, "Delete url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx = logger.WithUsername(ctx, username)

		if err := st.DeleteURL(ctx, username, alias); err != nil {
			err := fmt.Errorf("deleting url from storage: %w", err)
			slog.ErrorContext(ctx, "Delete url: "+err.Error())
			if errors.Is(err, storage.ErrAliasNotFound) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cache.DeleteURL(ctx, alias); err != nil {
			slog.WarnContext(ctx, "Delete url: deleting url from cache: "+err.Error())
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.DebugContext(ctx, "Deleted url")
	}
}
