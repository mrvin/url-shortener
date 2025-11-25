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

type URLDeleter interface {
	DeleteURL(ctx context.Context, username, alias string) error
}

func NewDeleteURL(deleter URLDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		username, err := logger.GetUsernameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Delete url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := deleter.DeleteURL(req.Context(), username, alias); err != nil {
			err := fmt.Errorf("deleting url from storage: %w", err)
			slog.ErrorContext(req.Context(), "Delete url: "+err.Error(), slog.String("alias", alias))
			if errors.Is(err, storage.ErrAliasNotFound) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "Deleted url", slog.String("alias", alias))
	}
}
