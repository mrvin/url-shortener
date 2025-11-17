package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

type URLsGetter interface {
	GetURLs(ctx context.Context, username string, limit, offset uint64) ([]storage.URL, uint64, error)
}

type ResponseGetURLs struct {
	URLs   []storage.URL `json:"urls"`
	Total  uint64        `json:"total"`
	Status string        `json:"status"`
}

func NewGetURLs(getter URLsGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error

		limit := uint64(defaultLimit)
		limitStr := req.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err = strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect limit value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		offset := uint64(defaultOffset)
		offsetStr := req.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err = strconv.ParseUint(offsetStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect offset value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}

		username, err := logger.GetUsernameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Get urls: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		urls, total, err := getter.GetURLs(req.Context(), username, limit, offset)
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
			slog.Uint64("total", total),
		)
	}
}
