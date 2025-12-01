package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
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

func NewGetURLs(getter URLsGetter) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		var err error
		ctx := req.Context()

		limit := uint64(defaultLimit)
		limitStr := req.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err = strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				return ctx, http.StatusBadRequest, fmt.Errorf("incorrect limit value: %w", err)
			}
		}
		offset := uint64(defaultOffset)
		offsetStr := req.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err = strconv.ParseUint(offsetStr, 10, 64)
			if err != nil {
				return ctx, http.StatusBadRequest, fmt.Errorf("incorrect offset value: %w", err)
			}
		}

		username, err := logger.GetUsernameFromCtx(ctx)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("get user name from ctx: %w", err)
		}

		urls, total, err := getter.GetURLs(ctx, username, limit, offset)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("getting urls from storage: %w", err)
		}

		// Write json response
		response := ResponseGetURLs{
			URLs:   urls,
			Total:  total,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("marshal response: %w", err)
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := res.Write(jsonResponse); err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("write response: %w", err)
		}

		return ctx, http.StatusOK, nil
	}
}
