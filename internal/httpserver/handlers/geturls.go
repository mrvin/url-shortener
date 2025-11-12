package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type URLsGetter interface {
	GetURLs(ctx context.Context, username string) ([]storage.URL, int64, error)
}

type ResponseGetURLs struct {
	URLs   []storage.URL `json:"urls"`
	Total  int64         `json:"total"`
	Status string        `json:"status"`
}

func NewGetURLs(getter URLsGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		username, err := logger.GetUsernameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Get urls: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		urls, total, err := getter.GetURLs(req.Context(), username)
		if err != nil {
			slog.ErrorContext(req.Context(), "Get urls: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseGetURLs{
			URLs:   urls,
			Total:  total,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Get urls: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Get urls: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(req.Context(), "Get urls",
			slog.Int64("total", total),
		)
	}
}
