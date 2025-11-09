package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, userName, alias string) error
}

func NewDeleteURL(deleter URLDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		userName, err := logger.GetUserNameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Delete url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := deleter.DeleteURL(req.Context(), userName, alias); err != nil {
			err := fmt.Errorf("failed delete url: %w", err)
			slog.ErrorContext(req.Context(), "Delete url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "Deleted url", slog.String("alias", alias))
	}
}
