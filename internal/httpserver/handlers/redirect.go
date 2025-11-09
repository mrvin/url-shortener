package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func NewRedirect(getter URLGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		resURL, err := getter.GetURL(req.Context(), alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				err := fmt.Errorf("url not found: %w", err)
				slog.ErrorContext(req.Context(), "Redirect: "+err.Error(), slog.String("alias", alias))
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			err := fmt.Errorf("failed get url: %w", err)
			slog.ErrorContext(req.Context(), "Redirect: "+err.Error(), slog.String("alias", alias))
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to found url
		http.Redirect(res, req, resURL, http.StatusFound)

		slog.InfoContext(req.Context(), "Redirect",
			slog.String("alias", alias),
			slog.String("url", resURL),
		)
	}
}
