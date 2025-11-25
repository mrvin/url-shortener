package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
)

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func NewRedirect(getter URLGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		url, err := getter.GetURL(req.Context(), alias)
		if err != nil {
			err = fmt.Errorf("getting url from storage: %w", err)
			if errors.Is(err, storage.ErrAliasNotFound) {
				slog.ErrorContext(req.Context(), "Redirect: "+err.Error(), slog.String("alias", alias))
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			slog.ErrorContext(req.Context(), "Redirect: "+err.Error(), slog.String("alias", alias))
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to found url
		http.Redirect(res, req, url, http.StatusFound)

		slog.InfoContext(req.Context(), "Redirect",
			slog.String("alias", alias),
			slog.String("url", url),
		)
	}
}
