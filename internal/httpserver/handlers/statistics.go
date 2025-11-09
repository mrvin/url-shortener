package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type CountGetter interface {
	GetCountURL(ctx context.Context, userName, alias string) (int64, error)
}

type ResponseGetCount struct {
	Count  int64  `json:"count"`
	Status string `json:"status"`
}

func NewGetCount(getter CountGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		userName, err := logger.GetUserNameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Get count: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := getter.GetCountURL(req.Context(), userName, alias)
		if err != nil {
			err := fmt.Errorf("get count: %w", err)
			slog.InfoContext(req.Context(), "Get count: "+err.Error(), slog.String("alias", alias))
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Write json response
		response := ResponseGetCount{
			Count:  count,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Get count: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Get count: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(req.Context(), "Get count",
			slog.Int64("count", count),
			slog.String("alias", alias),
		)
	}
}
