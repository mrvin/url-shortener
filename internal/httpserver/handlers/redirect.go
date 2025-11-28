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

type DBURLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
	CountIncrement(alias string) error
}

type CacheURLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
	SetURL(ctx context.Context, url, alias string) error
}

func NewRedirect(st DBURLGetter, cache CacheURLGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")
		ctx := logger.WithAlias(req.Context(), alias)

		url, err := cache.GetURL(ctx, alias)
		if err != nil {
			slog.WarnContext(ctx, "Redirect: getting url from cache: "+err.Error())
		}
		if url == "" {
			url, err = st.GetURL(ctx, alias)
			if err != nil {
				err = fmt.Errorf("getting url from storage: %w", err)
				if errors.Is(err, storage.ErrAliasNotFound) {
					slog.ErrorContext(ctx, "Redirect: "+err.Error())
					httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
					return
				}
				slog.ErrorContext(ctx, "Redirect: "+err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := cache.SetURL(ctx, url, alias); err != nil {
				slog.WarnContext(ctx, "Redirect: "+err.Error())
			}
		} else {
			go func() {
				if err := st.CountIncrement(alias); err != nil {
					slog.WarnContext(ctx, "Redirect: "+err.Error())
				}
			}()
		}

		// redirect to found url
		http.Redirect(res, req, url, http.StatusFound)

		slog.DebugContext(ctx, "Redirect",
			slog.String("url", url),
		)
	}
}
